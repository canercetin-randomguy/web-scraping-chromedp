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

func SignInHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signin.html",
		gin.H{
			"CallbackURL": SigninCallbackURL,
			"SignupURL":   SignupURL,
			"SigninURL":   SigninURL,
		})
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
		if err != nil {
			loggingUtil.Error("Error while connecting to database.", zap.Error(err))
		}
		err, hashedPassword := dbConnection.GetExistingUserPassword(LoginJSON.Username)
		if err != nil {
			loggingUtil.Error(fmt.Sprintf("Error while getting password of the user %s", LoginJSON.Username), zap.Error(err))
		}
		err = CompareHash(hashedPassword, LoginJSON)
		if err != nil {
			loggingUtil.Error(fmt.Sprintf("User %s entered the password wrong.", LoginJSON.Username), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status": "failed",
			})
		} else {
			// Then set an auth token cookie.
			// get a random auth token first.
			auth := sqlpkg.RandStringBytesMaskImprSrcSB(60)
			// nuke the auth token to the clients username.
			err = dbConnection.InsertAuthenticationToken(LoginJSON.Username, auth)
			if err != nil {
				log.Println(err)
			}
			if err != nil {
				log.Println(err)
			}
			c.JSON(http.StatusOK, gin.H{
				"status":                 "success",
				"usernameCookie":         LoginJSON.Username,
				"usernameCookieExpires":  3600,
				"authTokenCookie":        auth,
				"authTokenCookieExpires": 3600,
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
