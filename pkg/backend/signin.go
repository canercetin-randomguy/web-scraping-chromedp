package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignInHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signin.html",
		gin.H{
			"CallbackURL": SigninCallbackURL,
		})
}

func SigninFormJSONBinding(c *gin.Context) {
	var LoginJSON = SignInFormBinding{}
	// Bind the json to the user credentials struct.
	err := c.BindJSON(&LoginJSON)
	if err != nil {
		fmt.Println(err)
	}
	// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
	// TODO: get password from username and compare the hash with plain text password.
}

func CompareHash(pwd []byte, nuserInfo SignInFormBinding) error {
	return nil
}
