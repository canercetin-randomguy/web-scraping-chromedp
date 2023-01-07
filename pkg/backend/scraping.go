package backend

import (
	"canercetin/pkg/links"
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
		if !strings.Contains(origin, "localhost") {
			loggingUtil.Infow("Someone tried to access the endpoint from outside localhost.",
				"utility", "ScrapingFormJSONBinding")
			c.Status(http.StatusForbidden)
			return
		}
		var ScrapingJSON = ScrapingFormBinding{}
		// Bind the json to the scraping struct.
		err := c.BindJSON(&ScrapingJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
		}
		maxDepthInteger, err := strconv.Atoi(ScrapingJSON.MaxDepth)
		if err != nil {
			loggingUtil.Errorw("Error while converting max depth to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
		}
		linkLimitInteger, err := strconv.Atoi(ScrapingJSON.LinkLimit)
		if err != nil {
			loggingUtil.Errorw("Error while converting link limit to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
		}
		linkJson, fileNum := links.FindLinks(ScrapingJSON.MainLink, maxDepthInteger, ScrapingJSON.Username, linkLimitInteger)
		clientFolderpath := fmt.Sprintf("./logs/%s", ScrapingJSON.Username)
		jsonFilepath := fmt.Sprintf("%s/result_%s_%s_%d.json", clientFolderpath, ScrapingJSON.Username, time.Now().Format("20060102"), fileNum)
		// save the linkJson to a file
		err = ioutil.WriteFile(jsonFilepath, []byte(linkJson), 0644)
		if err != nil {
			loggingUtil.Errorw("Error while writing to JSON file", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"jsonFilepath", jsonFilepath)
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
		}
		data, err := ioutil.ReadFile(csvFilepath)
		if err != nil {
			loggingUtil.Errorw("Error while reading CSV file.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"csvFilepath", csvFilepath)
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			// send the json back
			"json": linkJson,
			"csv":  string(data),
		})
	}
}
