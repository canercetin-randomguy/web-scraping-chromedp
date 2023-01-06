package backend

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

// SignupFormJSONBinding sets JSON data that has arrived from signup.html's fetch request.
func SignupFormJSONBinding(loggingUtil *zap.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// close the endpoint from anyone but localhost, so signup.html can send a POST request but no one else.
		origin := c.Request.Header.Get("Origin")
		if !strings.Contains(origin, "localhost") {
			loggingUtil.Error("Someone tried to access the endpoint from outside of localhost.")
			c.Status(http.StatusForbidden)
			return
		}
		var LoginJSON = SignUpFormBinding{}
		// Bind the json to the user credentials struct.
		err := c.BindJSON(&LoginJSON)
		if err != nil {
			loggingUtil.Error("Error while binding signup form json to struct", zap.Error(err))
			return
		}
		// if the username or password is empty, return failed json.
		if LoginJSON.Username == "" || LoginJSON.Password == "" || LoginJSON.Email == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": "failed",
				"error":  "one of the fields are empty",
			})
			loggingUtil.Error("User tried to sign up with empty fields", zap.String("username", LoginJSON.Username))
			c.Status(http.StatusBadRequest)
			return
		}
		// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
		err = hashAndSalt([]byte(LoginJSON.Password), 16, LoginJSON)
		if err != nil {
			loggingUtil.Error(fmt.Sprintf("Error while registering user %s to database", LoginJSON.Username), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
		} else {
			// Send a response to the client that the user has been created.
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
			})
		}
	})
}

func hashAndSalt(pwd []byte, minCost int, userInfo SignUpFormBinding) error {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, minCost)
	if err != nil {
		return err
	}
	dbConnection := sqlpkg.SqlConn{}
	err = dbConnection.GetSQLConn("clients")
	if err != nil {
		return err
	}
	err = dbConnection.CreateNewUser(userInfo.Username, hash, userInfo.Email, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	err = dbConnection.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// SignupPage is a literally fleshed out signup page just consisting three input fields with a submit button.
func SignupPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{
			// This will be ponged back to server when client clicks the submit button.
			"CallbackURL": SignupCallbackURL,
			"SignupURL":   SignupURL,
			"SigninURL":   SigninURL,
		},
	)
}
