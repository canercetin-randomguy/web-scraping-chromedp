package backend

import (
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// SignupFormJSONBinding sets JSON data that has arrived from signup.html's fetch request.
func SignupFormJSONBinding(c *gin.Context) {
	var LoginJSON = LoginCredentials{}
	// Bind the josn to the user credentials struct.
	err := c.BindJSON(&LoginJSON)
	if err != nil {
		fmt.Println(err)
	}
	// Hash the password and salt it with 16 min cost, this can change. Then create a new user with the LoginJSON struct.
	err = hashAndSalt([]byte(LoginJSON.Password), 16, LoginJSON)
	if err != nil {
		// Send a response to the client that the user already exists.
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}
func hashAndSalt(pwd []byte, minCost int, userInfo LoginCredentials) error {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	_, err := bcrypt.GenerateFromPassword(pwd, minCost)
	if err != nil {
		log.Println(err)
	}
	// TODO: Nuke the userInfo to database.
	return nil
}

// SignupPage is a literally fleshed out signup page just consisting three input fields with a submit button.
func SignupPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{
			"CallbackURL": "callback",
		},
	)
}

// StartWebPageBackend starts the web page backend. Literally.
func StartWebPageBackend(localPort int) {
	r := gin.Default()
	// This is used for hiding printing one hundred of lines of loading static files.
	// If you want to see which files are loaded you can remove this line.
	gin.SetMode(gin.ReleaseMode)
	r.GET("/signup", SignupPage)
	r.POST("/callback", SignupFormJSONBinding)
	r.GET("/exists", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"exists.html",
			gin.H{},
		)
	})
	r.GET("/success", func(c *gin.Context) {
		c.Writer.Write([]byte("Success!"))
	})
	r.HTMLRender = ginview.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Run(fmt.Sprintf(":%d", localPort))
	// r.Run(":6969")
}
