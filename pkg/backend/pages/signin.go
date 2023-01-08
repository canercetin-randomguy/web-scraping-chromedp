package pages

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func SignInHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := c.Cookie("username")
		if err != nil {
			loggingUtil.Errorw("Error while getting username cookie", zap.Error(err),
				"username", username)
			c.Status(http.StatusInternalServerError)
			return
		}
		authCookie, err := c.Cookie("authtoken")
		if err != nil {
			loggingUtil.Errorw("Error while getting auth token cookie", zap.Error(err),
				"username", username)
			c.Status(http.StatusInternalServerError)
			return
		}
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Errorw("Error while opening database connection", zap.Error(err),
				"username", username)
			c.Status(http.StatusInternalServerError)
			return
		}
		defer dbConnection.DB.Close()
		auth, err := dbConnection.RetrieveAuthenticationToken(username)
		if err != nil {
			loggingUtil.Errorw("Error while retrieving authentication token", zap.Error(err),
				"username", username)
			c.Status(http.StatusInternalServerError)
			return
		}
		if authCookie != auth {
			c.HTML(
				http.StatusOK,
				"signin.html",
				gin.H{
					"CallbackURL": SigninCallbackURL,
					"SignupURL":   SignupURL,
					"SigninURL":   SigninURL,
				})
		} else {
			c.Redirect(http.StatusMovedPermanently, HomePath)
		}
	}
}

// SigninFormJSONBinding sets JSON data that has arrived from signin.html's fetch request.
func SigninFormJSONBinding(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var LoginJSON = SignInFormBinding{}
		// Bind the json to the user credentials struct.
		err := c.BindJSON(&LoginJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "SigninFormJSONBinding")
		}
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		defer dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Error("Error while connecting to database.", zap.Error(err))
		}
		err, hashedPassword := dbConnection.GetExistingUserPassword(LoginJSON.Username)
		if err != nil {
			loggingUtil.Error(fmt.Sprintf("Error while getting password of the user %s", LoginJSON.Username), zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "Error while getting password of the user",
			})
			c.Status(http.StatusBadGateway)
			return
		}
		err = CompareHash(hashedPassword, LoginJSON)
		if err != nil {
			loggingUtil.Error(fmt.Sprintf("User %s entered the password wrong.", LoginJSON.Username), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status": "failed",
			})
			c.Status(http.StatusBadGateway)
			return
		} else {
			// Then set an auth token cookie.
			// get a random auth token first.
			auth := sqlpkg.RandStringBytesMaskImprSrcSB(60)
			// nuke the auth token to the clients username.
			err = dbConnection.InsertAuthenticationToken(LoginJSON.Username, auth)
			if err != nil {
				log.Println(err)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":                 "success",
				"usernameCookie":         LoginJSON.Username,
				"usernameCookieExpires":  1,
				"authTokenCookie":        auth,
				"authTokenCookieExpires": 1,
			})
		}
		err = dbConnection.CloseConn()
		if err != nil {
			log.Println(err)
		}
	}
}

func CompareHash(pwd []byte, userInfo SignInFormBinding) error {
	err := bcrypt.CompareHashAndPassword(pwd, []byte(userInfo.Password))
	if err != nil {
		return err
	}
	return nil
}
