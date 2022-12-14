package pages

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func RestrictCallbackAccess(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		if !strings.Contains(ipAddress, "127.0.0.1") {
			loggingUtil.Infow("User tried to access the endpoint from outside localhost.",
				"utility", "RestrictCallbackAccess")
			c.Status(http.StatusForbidden)
			return
		}
	}
}

func RestrictPageAccess(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
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
			dbConnection := sqlpkg.SqlConn{}
			err := dbConnection.GetSQLConn("clients")
			defer dbConnection.DB.Close()
			if err != nil {
				loggingUtil.Info(fmt.Sprintf("Could not open database connection while handling user login %s.", user), zap.Error(err),
					"utility", "RestrictPageAccess",
					"client", user)
				c.Status(http.StatusInternalServerError)
				return
			}
			// then search for auth token in DB.
			auth, err := dbConnection.RetrieveAuthenticationToken(user)
			if err != nil {
				loggingUtil.Info(fmt.Sprintf("User %s authentication token could not retrieved from database.", user), zap.Error(err),
					"utility", "RestrictPageAccess",
					"client", user)
				c.Status(http.StatusInternalServerError)
				return
			}
			// If the request is not coming from localhost, a client may be trying to access the endpoint to download the file.
			// Or they may be trying to sign up or sign in.
			// Check if they want to download a file.
			if strings.Contains(c.Request.URL.String(), "storage") {
				// Check if the user is the owner of the file.
				if !strings.Contains(c.Request.URL.String(), user) {
					loggingUtil.Infow("Someone tried to access the endpoint from outside localhost.",
						"utility", "RestrictPageAccess")
					c.Status(http.StatusForbidden)
					c.Redirect(302, HomePath)
					return
				} else {
					authCookie, err := c.Cookie("authtoken")
					if err != nil {
						loggingUtil.Infow("User cookie cannot be retrieved.",
							"utility", "RestrictPageAccess",
							"client", "user")
						c.Status(http.StatusForbidden)
						c.Redirect(302, HomePath)
						return
					}
					if authCookie != auth {
						loggingUtil.Infow("Someone tried to access the endpoint from outside client.",
							"utility", "RestrictPageAccess",
							"client", user)
						c.Status(http.StatusForbidden)
						c.Redirect(302, HomePath)
						return
					}
				}
			}
			// if auth token is not found, redirect to login page, if we are not in the login or signup page.
			if auth == "" && (strings.Contains(c.Request.URL.Path, "signin") || strings.Contains(c.Request.URL.Path, "signup")) {
				// then delete the cookie.
				c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
				c.SetCookie("username", "", -1, "/", "localhost", false, true)
				loggingUtil.Info(fmt.Sprintf("User %s tried to access home page without having an auth token", user), zap.Error(err),
					"utility", "RestrictPageAccess",
					"client", user)
				c.Redirect(302, SigninPath)
				return
			}
			// if auth token is found, check if it is valid.
			authCookie, err := c.Cookie("authtoken")
			if auth != authCookie && (!strings.Contains(c.Request.URL.Path, "signin") || !strings.Contains(c.Request.URL.Path, "signup")) {
				// then delete the cookie.
				c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
				c.SetCookie("username", "", -1, "/", "localhost", false, true)
				loggingUtil.Infow(fmt.Sprintf("User %s tried to access home page with an invalid auth token", user), zap.Error(err),
					"utility", "RestrictPageAccess",
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
