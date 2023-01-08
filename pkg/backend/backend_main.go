package backend

import (
	"canercetin/pkg/backend/pages"
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
	v1 := r.Group("/v1", pages.RestrictPageAccess(loggingUtil))
	v1.StaticFS("/storage", http.Dir("./results/staticfs"))
	v1.GET("/signup", pages.SignupPage)
	v1.GET("/signin", pages.SignInHandler)
	v1.GET("/home", pages.HomeHandler(loggingUtil))
	v1.GET("/download", pages.DownloadPage(loggingUtil))
	v1.GET("/logout", pages.LogoutHandler(loggingUtil))
	// disallow any user except localhost to access /public endpoint.
	loggingUtil.Info("Starting backend on port " + fmt.Sprint(localPort))
	err := r.Run(fmt.Sprintf(":%d", localPort))
	if err != nil {
		return err
	}
	// r.Run(":6969")
	return nil
}
