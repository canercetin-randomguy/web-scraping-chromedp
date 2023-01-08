package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func LogoutHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the current domain
		domain := c.Request.Host
		fmt.Println(domain)
		c.SetCookie("authtoken", "", 0, "/v1", domain, false, false)
		c.SetCookie("username", "", 0, "/v1", domain, false, false)
		c.Redirect(http.StatusFound, SigninPath)
		return
	}
}
