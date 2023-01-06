package backend

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HomeHandler(loggingUtil *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Search for username in cookies.
		user, _ := c.Cookie("username")
		dbConnection := sqlpkg.SqlConn{}
		err := dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("Could not open database connection while handling user login %s.", user), zap.Error(err))
		}
		// then search for auth token in DB.
		auth, err := dbConnection.RetrieveAuthenticationToken(user)
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("User %s authentication token could not retrieved from database.", user), zap.Error(err))
		}
		// if auth token is not found, redirect to login page.
		if auth == "" {
			// then delete the cookie.
			c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
			c.SetCookie("username", "", -1, "/", "localhost", false, true)
			c.Redirect(302, "/signin")
			loggingUtil.Info(fmt.Sprintf("User %s tried to access home page without having an auth token", user), zap.Error(err))
			return
		}
		// if auth token is found, compare it with the cookie.
		cookie, _ := c.Cookie("authtoken")
		if cookie != auth { // if they don't match, redirect to login page.
			c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
			c.SetCookie("username", "", -1, "/", "localhost", false, true)
			c.Redirect(302, "/signin")
			loggingUtil.Info(fmt.Sprintf("User %s auth token is not matching the one with the one in database.", user), zap.Error(err))
			// then delete the cookie.
			return
		}
		// query to see if user is on free plan.
		planName, err := dbConnection.RetrieveUserPackageDetails(user)
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("User %s package details couldnt be retrieved.", user), zap.Error(err))
			return
		}
		limitAmount, err := dbConnection.RetrieveUserLinkLimit(user)
		if err != nil {
			loggingUtil.Info(fmt.Sprintf("User %s package limit details couldnt be retrieved.", user), zap.Error(err))
			return
		}
		err = dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Info("Couldn't close the database connection.", zap.Error(err))
			return
		}
		c.HTML(
			200,
			"home.html",
			gin.H{
				// When users submit the form, it will be sent to /scrape as a post request. Dont forget to hide it.
				"CallbackURL": ScrapingURL,
				"Username":    user,
				"Plan":        planName,
				"Limit":       limitAmount,
			},
		)
	})
}
