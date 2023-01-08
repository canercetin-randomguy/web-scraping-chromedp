package pages

import (
	"canercetin/pkg/backend"
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HomeHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Search for username in cookies.
		user, _ := c.Cookie("username")
		dbConnection := sqlpkg.SqlConn{}
		err := dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("Could not open database connection while handling user login %s.", user), zap.Error(err))
		}
		// query to see if user is on free plan.
		planName, err := dbConnection.RetrieveUserPackageDetails(user)
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("User %s package details couldnt be retrieved.", user), zap.Error(err))
			return
		}
		// query to see link limit.
		limitAmount, err := dbConnection.RetrieveUserLinkLimit(user)
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("User %s package limit details couldnt be retrieved.", user), zap.Error(err))
			return
		}
		// close the database connection.
		err = dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Info("Couldn't close the database connection.", zap.Error(err))
			return
		}
		// serve the html.
		c.HTML(
			200,
			"home.html",
			gin.H{
				// When users submit the form, it will be sent to /scrape as a post request. Dont forget to hide it.
				"CallbackURL": backend.ScrapingURL,
				"Username":    user,
				"Plan":        planName,
				"Limit":       limitAmount,
			},
		)
	}
}
