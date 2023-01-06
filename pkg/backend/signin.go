package backend

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
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
func SigninFormJSONBinding(c *gin.Context) {
	// close the endpoint from anyone but localhost, so signin.html can send a POST request but no one else.
	origin := c.Request.Header.Get("Origin")
	if !strings.Contains(origin, "localhost") {
		c.Status(http.StatusForbidden)
		return
	}
	var LoginJSON = SignInFormBinding{}
	// Bind the json to the user credentials struct.
	err := c.BindJSON(&LoginJSON)
	if err != nil {
		fmt.Println(err)
	}
	// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
	// TODO: get password from username and compare the hash with plain text password.
	dbConnection := sqlpkg.SqlConn{}
	err = dbConnection.GetSQLConn("clients")
	if err != nil {
		log.Println(err)
	}
	err, hashedPassword := dbConnection.GetExistingUserPassword(LoginJSON.Username)
	if err != nil {
		log.Println(err)
	}
	err = CompareHash(hashedPassword, LoginJSON)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "failed",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func CompareHash(pwd []byte, userInfo SignInFormBinding) error {
	err := bcrypt.CompareHashAndPassword(pwd, []byte(userInfo.Password))
	if err != nil {
		return err
	}
	return nil
}
