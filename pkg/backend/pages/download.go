package pages

import (
	"canercetin/pkg/sqlpkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DownloadPage(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the username from the cookie.
		user, _ := c.Cookie("username")
		dbConnection := sqlpkg.SqlConn{}
		err := dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Info("Could not open database connection while handling user login.", zap.Error(err),
				"utility", "DownloadPage")
		}
		str, err := dbConnection.RetrieveFileLinks(user)
		if err != nil {
			loggingUtil.Errorw("Could not retrieve file links from database.", zap.Error(err),
				"utility", "DownloadPage",
				"client", user)
		}
		err = dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Errorw("Could not close database connection.", zap.Error(err),
				"utility", "DownloadPage",
				"client", user)
		}
		c.HTML(200, "download.html", gin.H{
			// pass the testStruct to the template.
			"teststruct": str,
		})

	}
}
