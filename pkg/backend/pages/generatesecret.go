package pages

import (
	"canercetin/pkg/sqlpkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SecretPageHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "secretkey.html", gin.H{
			"SecretKeyCallbackPath": SecretKeyCallbackURL,
		})
	}
}
func SecretKeyFormJSONBinding(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var secretKeyForm SecretKeyFormBinding
		err := c.ShouldBindJSON(&secretKeyForm)
		if err != nil {
			loggingUtil.Errorw("Error while binding JSON", "error", err)
			c.JSON(400, gin.H{
				"error": "Error while binding JSON",
			})
			return
		}
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		if err != nil {
			loggingUtil.Errorw("Error while opening database connection", zap.Error(err))
			c.Status(500)
			return
		}
		defer dbConnection.DB.Close()
		auth, err := dbConnection.RetrieveAuthenticationToken(secretKeyForm.Username)
		if auth != secretKeyForm.AuthKey {
			loggingUtil.Errorw("Authentication key does not match", zap.Error(err))
			c.JSON(400, gin.H{
				"error": "Auth key is not valid",
			})
			return
		}
		secretKey := sqlpkg.RandStringBytesMaskImprSrcSB(60)
		err = dbConnection.InsertSecretKey(secretKeyForm.Username, secretKey)
		if err != nil {
			loggingUtil.Errorw("Error while inserting secret key", zap.Error(err))
			c.Status(500)
			return
		}
		c.JSON(200, gin.H{
			"status":    "success",
			"secretKey": secretKey,
		})
	}
}
