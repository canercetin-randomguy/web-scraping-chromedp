package api

import (
	"bytes"
	"canercetin/pkg/backend/pages"
	"canercetin/pkg/sqlpkg"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func CrawlHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var CrawlJSON = CrawlPOSTBinding{}
		err := c.BindJSON(&CrawlJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Error while binding JSON to struct, please make sure you have sent raw JSON data.",
			})
		}
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		defer dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Errorw("Error while connecting to database.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while connecting to database. Caner probably screwed something up.",
			})
			return
		}
		auth, err := dbConnection.RetrieveAuthenticationToken(CrawlJSON.Username)
		if err != nil {
			loggingUtil.Errorw("Error while retrieving authentication token from database.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while retrieving authentication token from database. Contact the developer.",
			})
			return
		}
		if auth != CrawlJSON.AuthKey {
			loggingUtil.Errorw("Authentication token does not match.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Authentication token does not match.",
			})
			return
		}
		linkLimit, err := dbConnection.RetrieveUserLinkLimit(CrawlJSON.Username)
		if err != nil {
			loggingUtil.Errorw("Error while retrieving link limit from database.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while retrieving link limit from database. Contact the developer.",
			})
			return
		}
		linkLimitString := strconv.Itoa(linkLimit)
		if linkLimitString == "0" {
			loggingUtil.Errorw("User has reached the link limit.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "failed",
				"message": "User has reached the link limit.",
			})
			return
		}
		// make a request to ScrapingLink
		jsonBody := pages.ScrapingFormBinding{
			Username:  CrawlJSON.Username,
			AuthKey:   CrawlJSON.AuthKey,
			MaxDepth:  CrawlJSON.MaxDepth,
			MainLink:  CrawlJSON.MainLink,
			LinkLimit: linkLimitString,
		}
		// marshal the jsonBody
		jsonBodyBytes, err := json.Marshal(jsonBody)
		if err != nil {
			loggingUtil.Errorw("Error while marshaling JSON.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while marshaling JSON. Contact the developer.",
			})
			return
		}
		req, err := http.NewRequest(http.MethodPost, pages.ScrapingPath, bytes.NewReader(jsonBodyBytes))
		if err != nil {
			loggingUtil.Errorw("Error while making a request to scraping link.", zap.Error(err),
				"utility", "CrawlHandler",
				"link", pages.ScrapingPath)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			loggingUtil.Errorw("Error while reading response body.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while reading response body. Contact the developer.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"response": responseData,
		})
	}
}
