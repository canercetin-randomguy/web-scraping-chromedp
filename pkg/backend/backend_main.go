package backend

import (
	"canercetin/pkg/logger"
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StartWebPageBackend starts the web page backend. Literally.
func StartWebPageBackend(localPort int) error {
	backendLogFilepath := logger.CreateNewFile("./logs/backend")
	loggingUtil, err2 := logger.NewLoggerWithFile(backendLogFilepath)
	if err2 != nil {
		return err2
	}
	r := gin.Default()
	r.HTMLRender = ginview.Default()
	r.LoadHTMLGlob("templates/*.html")
	v1 := r.Group("/v1", RestrictSysAccess(loggingUtil))
	v1.StaticFS("/storage", http.Dir("./results/staticfs"))
	v1.GET("/signup", SignupPage)
	v1.GET("/signin", SignInHandler)
	v1.GET("/home", HomeHandler(loggingUtil))
	v1.GET("/download", DownloadPage(loggingUtil))
	v1.GET("/logout", LogoutHandler(loggingUtil))
	// If client hits submit button, make a post request to this endpoint and this endpoint will return a json. T
	v1.POST("/signin/callback", SigninFormJSONBinding(loggingUtil))
	v1.POST("/home/scraping/callback", ScrapingFormJSONBinding(loggingUtil))
	v1.POST("/signup/callback", SignupFormJSONBinding(loggingUtil))
	v1.Static("/static", "./templates/static/")
	// disallow any user except localhost to access /public endpoint.
	loggingUtil.Info("Starting backend on port " + fmt.Sprint(localPort))
	err := r.Run(fmt.Sprintf(":%d", localPort))
	if err != nil {
		return err
	}
	// r.Run(":6969")
	return nil
}
