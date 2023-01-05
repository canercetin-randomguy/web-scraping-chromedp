package backend

import (
	"canercetin/pkg/sqlpkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// SignupFormJSONBinding sets JSON data that has arrived from signup.html's fetch request.
func SignupFormJSONBinding(c *gin.Context) {
	var LoginJSON = SignUpFormBinding{}
	// Bind the json to the user credentials struct.
	err := c.BindJSON(&LoginJSON)
	if err != nil {
		fmt.Println(err)
	}
	// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
	err = hashAndSalt([]byte(LoginJSON.Password), 16, LoginJSON)
	if err != nil {
		fmt.Println(err)
		// Send a response to the client that the user already exists.
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
		})
	} else {
		// Send a response to the client that the user has been created.
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

func hashAndSalt(pwd []byte, minCost int, userInfo SignUpFormBinding) error {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, minCost)
	if err != nil {
		log.Println(err)
	}
	// TODO: Nuke the userInfo to database.
	dbConnection := sqlpkg.SqlConn{}
	dbConnection.GetSQLConn("clients")
	dbConnection.CreateNewUser(userInfo.Username, hash, userInfo.Email, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// SignupPage is a literally fleshed out signup page just consisting three input fields with a submit button.
func SignupPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{
			// This will be ponged back to server when client clicks the submit button.
			"CallbackURL": SignupCallbackURL,
		},
	)
}
