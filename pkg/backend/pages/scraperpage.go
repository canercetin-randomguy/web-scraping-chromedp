package pages

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
		var ScrapingJSON = ScrapingFormBinding{}
		// Bind the json to the scraping struct.
		err := c.BindJSON(&ScrapingJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Error while binding JSON to struct, please make sure you have sent raw JSON data.",
			})
			return
		}
		maxDepthInteger, err := strconv.Atoi(ScrapingJSON.MaxDepth)
		if err != nil {
			loggingUtil.Errorw("Error while converting max depth to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Error while converting max depth to integer, please contact the developer.",
			})
			return
		}
		linkLimitInteger, err := strconv.Atoi(ScrapingJSON.LinkLimit)
		if err != nil {
			loggingUtil.Errorw("Error while converting link limit to integer.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Error while converting link limit to integer, please contact the developer.",
			})
			return
		}
		linkJson, fileNum := links.FindLinks(ScrapingJSON.MainLink, maxDepthInteger, ScrapingJSON.Username, linkLimitInteger)
		// Create a new folder for saving downloaded files.
		err = logger.CreateNewFolder(fmt.Sprintf("./results/staticfs/%s", ScrapingJSON.Username))
		if err != nil {
			loggingUtil.Errorw("Error while creating new folder for user.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while creating new folder for user, please contact the developer.",
			})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while writing to JSON file, please contact the developer.",
			})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while converting JSON to file, please contact the developer.",
			})
			return
		}

		// Send the jsonFilepath and csvFilepath to the database.
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		defer dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Errorw("Error  opening DB connection while sending json and csv", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"jsonFilepath", jsonFilepath,
				"csvFilepath", csvFilepath)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error  opening DB connection while sending json and csv, please contact the developer.",
			})
			return
		}
		jsonAbsoluteFilepath := FileStoragePath + strings.ReplaceAll(jsonFilepath, "./results/staticfs", "")
		csvAbsoluteFilepath := FileStoragePath + strings.ReplaceAll(csvFilepath, "./results/staticfs", "")
		err = dbConnection.InsertFileLink(ScrapingJSON.Username, jsonAbsoluteFilepath, time.Now().Format("2006-01-02 15:04:05"), "json", ScrapingJSON.MainLink)
		err = dbConnection.InsertFileLink(ScrapingJSON.Username, csvAbsoluteFilepath, time.Now().Format("2006-01-02 15:04:05"), "csv", ScrapingJSON.MainLink)

		if err != nil {
			loggingUtil.Errorw("Error while reading CSV file.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding",
				"client", ScrapingJSON.Username,
				"csvFilepath", csvFilepath)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while reading CSV file, please contact the developer.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":           "ok",
			"csvDownloadLink":  fmt.Sprintf("%s%s", FileStoragePath, strings.ReplaceAll(csvFilepath, "./results/staticfs", "")),
			"jsonDownloadLink": fmt.Sprintf("%s%s", FileStoragePath, strings.ReplaceAll(jsonFilepath, "./results/staticfs", "")),
		})
	}
}
