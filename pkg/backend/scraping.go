package backend

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
		var LoginJSON = ScrapingFormBinding{}
		// Bind the json to the scraping struct.
		err := c.BindJSON(&LoginJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "ScrapingFormJSONBinding")
		}
		// TODO: Do the scraping.
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			// send the json back
			"json": LoginJSON,
		})
	}
}
