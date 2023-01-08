package pages

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func SignInHandler(c *gin.Context) {
	auth, err := c.Cookie("authtoken")
	user, err := c.Cookie("username")
	if err != nil {
		c.Redirect(http.StatusFound, SigninPath)
		return
	}
	if user != "" {
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		if err != nil {
			log.Println(err)
		}
		authDB, err := dbConnection.RetrieveAuthenticationToken(user)
		if err != nil {
			log.Println(err)
		}
		err = dbConnection.CloseConn()
		if err != nil {
			log.Println(err)
		}
		if auth == authDB {
			c.Redirect(http.StatusFound, HomePath)
			return
		}
	}
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
		// close the endpoint from anyone but localhost, so signin.html can send a POST request but no one else.
		origin := c.Request.Header.Get("Origin")
		ipAddress := c.ClientIP()
		fmt.Println("Origin: ", origin)
		fmt.Println("IP Address: ", ipAddress)
		if !strings.Contains(origin, "localhost") || !strings.Contains(ipAddress, "::1") {
			loggingUtil.Infow("User tried to access the endpoint from outside localhost.",
				"utility", "SigninFormJSONBinding")
			c.Status(http.StatusForbidden)
			return
		}
		var LoginJSON = SignInFormBinding{}
		// Bind the json to the user credentials struct.
		err := c.BindJSON(&LoginJSON)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON to struct.", zap.Error(err),
				"utility", "SigninFormJSONBinding")
		}
		// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
		// TODO: get password from username and compare the hash with plain text password.
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
			// If user is successfully logged in, set a cookie of clients username.
			c.SetCookie("username", LoginJSON.Username, 3600, "/", "localhost", false, true)
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
			c.SetCookie("authtoken", auth, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
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