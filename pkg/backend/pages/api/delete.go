package api

import (
	"canercetin/pkg/sqlpkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

func DeleteHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ToDeleteInfo = DeletePOSTBinding{}
		err := c.ShouldBindJSON(&ToDeleteInfo)
		if err != nil {
			loggingUtil.Errorw("Error binding JSON.", zap.Error(err),
				"Utility", "DeleteHandler")
		}
		if err != nil {
			loggingUtil.Errorw("Error getting sql connection.", zap.Error(err),
				"client", ToDeleteInfo.Username,
				"Utility", "DeleteHandler")
			return
		}
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		defer func(DB *sql.DB) {
			err = DB.Close()
			if err != nil {
				panic(err)
			}
		}(dbConnection.DB)
		auth, err := dbConnection.RetrieveAuthenticationToken(ToDeleteInfo.Username)
		if err != nil {
			loggingUtil.Errorw("Error while retrieving authentication token from database.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "failed",
				"message": "Error while retrieving authentication token from database. Contact the developer.",
			})
			return
		}
		if auth != ToDeleteInfo.AuthKey {
			loggingUtil.Errorw("Authentication token does not match.", zap.Error(err),
				"utility", "CrawlHandler")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "Authentication token does not match.",
			})
			return
		}
		// this will be used to delete the file link from the database
		err = dbConnection.DeleteFileLink(ToDeleteInfo.Username, ToDeleteInfo.FilePath)
		if err != nil {
			loggingUtil.Errorw("Error deleting file link.", zap.Error(err))
		}
		// this will be used to delete the file from the server
		ToDeleteInfo.FilePath = strings.ReplaceAll(ToDeleteInfo.FilePath, "/v1/storage", "/results/staticfs")
		err = os.Remove(ToDeleteInfo.FilePath)
		if err != nil {
			loggingUtil.Errorw("Error deleting file.", zap.Error(err),
				"client", ToDeleteInfo.Username)
		}
	}
}
