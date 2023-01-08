package backend

import (
	"canercetin/pkg/backend/pages"
	"canercetin/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

// StartWebPageLoopback is used for endpoints that users are not allowed to go in.
func StartWebPageLoopback(port int) error {
	// Only users from localhost can access this endpoint, because we will run it on 127.0.0.1:3131
	backendLogFilepath := logger.CreateNewFile("./logs/backend")
	loggingUtil, err2 := logger.NewLoggerWithFile(backendLogFilepath)
	if err2 != nil {
		return err2
	}
	r := gin.Default()
	r.Use(CORSMiddleware())
	pv := r.Group("/private", pages.RestrictCallbackAccess(loggingUtil))
	pv.POST("/signin/callback", pages.SigninFormJSONBinding(loggingUtil))
	pv.POST("/home/scraping/callback", pages.ScrapingFormJSONBinding(loggingUtil))
	pv.POST("/signup/callback", pages.SignupFormJSONBinding(loggingUtil))
	pv.POST("/secretkey/callback", pages.SecretKeyFormJSONBinding(loggingUtil))
	err := r.Run(fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	return nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
