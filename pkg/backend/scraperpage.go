package backend

import (
	"canercetin/pkg/links"
	"canercetin/pkg/logger"
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ScraperHandler(c *gin.Context) {

}
func ScrapingFormJSONBinding(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dont allow anyone to access this endpoint that is not coming from localhost.
		origin := c.Request.Header.Get("Origin")
		ipAddress := c.ClientIP()
		if !strings.Contains(origin, "localhost") || !strings.Contains(ipAddress, "::1") {
			if !strings.Contains(ipAddress, "127.0.0.1") {
				loggingUtil.Infow("Someone tried to access the endpoint from outside localhost.",
					"utility", "SigninFormJSONBinding")
				c.Status(http.StatusForbidden)
				return
			}
		}
		var ScrapingJSON = ScrapingFormBinding{}
		// Bind the json to the scraping struct.
		err := c.BindJSON(&ScrapingJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			return
		}
		maxDepthInteger, err := strconv.Atoi(ScrapingJSON.MaxDepth)
		if err != nil {
			loggingUtil.Errorw("Error while converting max depth to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			return
		}
		linkLimitInteger, err := strconv.Atoi(ScrapingJSON.LinkLimit)
		if err != nil {
			loggingUtil.Errorw("Error while converting link limit to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			return
		}
		linkJson, fileNum := links.FindLinks(ScrapingJSON.MainLink, maxDepthInteger, ScrapingJSON.Username, linkLimitInteger)
		// Create a new folder for saving downloaded files.
		err = logger.CreateNewFolder(fmt.Sprintf("./results/staticfs/%s", ScrapingJSON.Username))
		if err != nil {
			loggingUtil.Errorw("Error while creating new folder for user.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			return
		}
		clientFolderpath := fmt.Sprintf("./results/staticfs/%s", ScrapingJSON.Username)
		jsonFilepath := fmt.Sprintf("%s/result_%s_%s_%d.json", clientFolderpath, ScrapingJSON.Username, time.Now().Format("20060102"), fileNum)
		// save the linkJson to a file
		err = ioutil.WriteFile(jsonFilepath, []byte(linkJson), 0644)
		if err != nil {
			loggingUtil.Errorw("Error while writing to JSON file", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"jsonFilepath", jsonFilepath)
			return
		}
		// Convert the saved json file to csv
		csvFilepath := fmt.Sprintf("%s/result_%s_%s_%d.csv", clientFolderpath, ScrapingJSON.Username, time.Now().Format("20060102"), fileNum)
		err = links.ConvertJSONToCSV(jsonFilepath, csvFilepath, loggingUtil, ScrapingJSON.Username)
		if err != nil {
			loggingUtil.Errorw("Error while converting JSON to file.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"jsonFilepath", jsonFilepath,
				"csvFilepath", csvFilepath)
			return
		}

		// Send the jsonFilepath and csvFilepath to the database.
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Errorw("Error  opening DB connection while sending json and csv", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"jsonFilepath", jsonFilepath,
				"csvFilepath", csvFilepath)
			return
		}
		err = dbConnection.InsertFileLink(ScrapingJSON.Username, jsonFilepath, time.Now().Format("2006-01-02 15:04:05"), "json")
		err = dbConnection.InsertFileLink(ScrapingJSON.Username, csvFilepath, time.Now().Format("2006-01-02 15:04:05"), "csv")

		if err != nil {
			loggingUtil.Errorw("Error while reading CSV file.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"csvFilepath", csvFilepath)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
