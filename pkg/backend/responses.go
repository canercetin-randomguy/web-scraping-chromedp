package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleSuccess is used when a sign in is successful. Will be used to redirect to the main page.
func HandleSuccess(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"success.html",
		gin.H{},
	)
}
