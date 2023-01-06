package backend

import (
	"canercetin/pkg/logger"
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// StartWebPageBackend starts the web page backend. Literally.
func StartWebPageBackend(localPort int) error {
	backendLogFilepath := logger.CreateNewFile("./logs/backend")
	loggingUtil, err2 := logger.NewLoggerWithFile(backendLogFilepath)
	if err2 != nil {
		log.Println(err2)
	}
	r := gin.Default()
	// This is used for hiding printing one hundred of lines of loading static files.
	// If you want to see which files are loaded you can remove this line.
	gin.SetMode(gin.ReleaseMode)
	// Used for taking sign up.
	r.GET("/signup", SignupPage)
	// Used for handling sign up requests.
	r.POST("/signup/callback", SignupFormJSONBinding(loggingUtil))
	// If it exists.
	r.GET("/exists", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"exists.html",
			gin.H{},
		)
	})
	// If client successfully signs up, yeet him to the sign-in page.
	r.GET("/signin", SignInHandler)
	// If client hits submit button, make a post request to this endpoint and this endpoint will return a json.
	r.POST("/signin/callback", SigninFormJSONBinding(loggingUtil))
	// If client successfully signs in, yeet him to the home page.
	r.GET("/home", HomeHandler(loggingUtil))
	r.HTMLRender = ginview.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./templates/static/")
	loggingUtil.Info("Starting backend on port " + fmt.Sprint(localPort))
	err := r.Run(fmt.Sprintf(":%d", localPort))
	if err != nil {
		return err
	}
	// r.Run(":6969")
	return nil
}
