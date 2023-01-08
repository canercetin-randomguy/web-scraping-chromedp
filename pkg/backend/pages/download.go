package pages

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
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

func RestrictSysAccess(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// IF ENDPOINT NEEDS TO BE CLOSED, SET THAT ENDPOINT IN BACKEND_LOOPBACK.
		// Check the cookies to which client.
		user, _ := c.Cookie("username")
		// Is user cookie empty, and is the endpoint signin or signup?
		if user == "" && c.Request.URL.Path != "/v1/signin" && c.Request.URL.Path != "/v1/signup" {
			c.Status(http.StatusUnauthorized)
			c.Redirect(http.StatusFound, SigninPath)
			return
		}
		// Are we in the signup or signin page?
		if strings.Contains(c.Request.URL.Path, "signup") || strings.Contains(c.Request.URL.Path, "signin") {
			// do nothing, let the user access the page.
		} else {
			// If the request is not coming from localhost, a client may be trying to access the endpoint to download the file.
			// Or they may be trying to sign up or sign in.
			// Check if they want to download a file.
			if !strings.Contains(c.Request.URL.String(), user) && strings.Contains(c.Request.URL.String(), "storage") {
				loggingUtil.Infow("Someone tried to access the endpoint from outside localhost.",
					"utility", "RestrictSysAccess")
				c.Status(http.StatusForbidden)
				c.Redirect(302, HomePath)
				return
			}
			// If yes get a database connection.
			dbConnection := sqlpkg.SqlConn{}
			err := dbConnection.GetSQLConn("clients")
			if err != nil {
				loggingUtil.Info(fmt.Sprintf("Could not open database connection while handling user login %s.", user), zap.Error(err),
					"utility", "RestrictSysAccess",
					"client", user)
			}
			// then search for auth token in DB.
			auth, err := dbConnection.RetrieveAuthenticationToken(user)
			if err != nil {
				loggingUtil.Info(fmt.Sprintf("User %s authentication token could not retrieved from database.", user), zap.Error(err),
					"utility", "RestrictSysAccess",
					"client", user)
			}
			// if auth token is not found, redirect to login page, if we are not in the login or signup page.
			if auth == "" && (strings.Contains(c.Request.URL.Path, "signin") || strings.Contains(c.Request.URL.Path, "signup")) {
				// then delete the cookie.
				c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
				c.SetCookie("username", "", -1, "/", "localhost", false, true)
				loggingUtil.Info(fmt.Sprintf("User %s tried to access home page without having an auth token", user), zap.Error(err),
					"utility", "RestrictSysAccess",
					"client", user)
				c.Redirect(302, SigninPath)
				return
			}
			// if auth token is found, check if it is valid.
			if auth != c.GetHeader("authtoken") && (strings.Contains(c.Request.URL.Path, "signin") || strings.Contains(c.Request.URL.Path, "signup")) {
				// then delete the cookie.
				c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
				c.SetCookie("username", "", -1, "/", "localhost", false, true)
				loggingUtil.Infow(fmt.Sprintf("User %s tried to access home page with an invalid auth token", user), zap.Error(err),
					"utility", "RestrictSysAccess",
					"client", user,
					"cookieAuthToken", auth,
					"headerAuthToken", c.GetHeader("authtoken"))
				c.Redirect(302, SigninPath)
				return
			}
			// if none of the above, let the client access the endpoint.
		}
	}
}
