package backend

import (
	"canercetin/pkg/links"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
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
		totalLinks, brokenLinks := links.FindLinks(ScrapingJSON.MainLink, maxDepthInteger, ScrapingJSON.Username, linkLimitInteger)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			// send the json back
			"json":        totalLinks,
			"brokenLinks": brokenLinks,
		})
	}
}
