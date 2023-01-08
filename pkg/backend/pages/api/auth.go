package api

import (
	"canercetin/pkg/sqlpkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandler(loggingUtil *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// print the form-data that has arrived.
		var AuthBinding = AuthPOSTBinding{}
		err := c.BindJSON(&AuthBinding)
		if err != nil {
			loggingUtil.Info("Could not bind json to AuthPOSTBinding struct.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Could not bind json to AuthPOSTBinding struct. Please make sure you are sending data raw JSON, if it doesnt work, please contact the developer.",
			})
			return
		}
		// get a new database connection
		dbConnection := sqlpkg.SqlConn{}
		err = dbConnection.GetSQLConn("clients")
		defer dbConnection.DB.Close()
		if err != nil {
			loggingUtil.Error("Error while connecting to database.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error while connecting to database. Caner probably screwed something up.",
			})
			return
		}
		err, dbUserpassword := dbConnection.GetExistingUserPassword(AuthBinding.Username)
		if err != nil {
			loggingUtil.Error("Error while getting user from database.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Error while getting user from database. Please make sure you are signed up.",
			})
			return
		}
		if string(dbUserpassword) == "" {
			loggingUtil.Error("User does not exist.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "User does not exist, please sign up. If you are already signed up, please contact the developer.",
			})
			return
		} else {
			if bcrypt.CompareHashAndPassword(dbUserpassword, []byte(AuthBinding.Password)) == nil {
				// do nothing
			} else {
				loggingUtil.Error("Wrong password.", zap.Error(err),
					zap.String("RequestIP", c.ClientIP()),
					zap.String("RequestURI", c.Request.RequestURI),
					zap.String("Client", AuthBinding.Username))
				c.JSON(400, gin.H{
					"status":  "failed",
					"message": "Wrong password. Please try again.",
				})
				return
			}
		}
		email, err := dbConnection.RetrieveEmail(AuthBinding.Username)
		if err != nil {
			loggingUtil.Error("Error while retrieving email.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error while retrieving email. Please contact the developer.",
			})
			return
		}
		if AuthBinding.Email != email {
			loggingUtil.Error("Wrong email.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Wrong email. Please try again.",
			})
			return
		}
		auth, err := dbConnection.RetrieveAuthenticationToken(AuthBinding.Username)
		if err != nil {
			loggingUtil.Error("Error while retrieving authentication token from database.", zap.Error(err),
				zap.String("RequestIP", c.ClientIP()),
				zap.String("RequestURI", c.Request.RequestURI),
				zap.String("Client", AuthBinding.Username))
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error while retrieving authentication token from database. Please contact the developer.",
			})
			return
		}
		if auth == "" || auth == " " {
			auth = sqlpkg.RandStringBytesMaskImprSrcSB(60)
			err = dbConnection.InsertAuthenticationToken(AuthBinding.Username, auth)
			if err != nil {
				loggingUtil.Error("Error while inserting authentication token to database.", zap.Error(err),
					zap.String("RequestIP", c.ClientIP()),
					zap.String("RequestURI", c.Request.RequestURI),
					zap.String("Client", AuthBinding.Username))
				c.JSON(500, gin.H{
					"status":  "failed",
					"message": "Error while inserting authentication token to database. Please contact the developer.",
				})
				return
			}
		}
		c.JSON(200, gin.H{
			"status":   "success",
			"username": AuthBinding.Username,
			"auth":     auth,
			"message":  "Successfully authenticated, please use the token for future requests.",
		})
	}
}
