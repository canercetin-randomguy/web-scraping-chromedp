package pages

import (
	"canercetin/pkg/backend"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func LogoutHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("authtoken", "", 0, "/", "localhost", false, true)
		c.SetCookie("username", "", 0, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, backend.SigninPath)
		return
	}
}
