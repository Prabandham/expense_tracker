package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
)

var redis = config.GetRedisConnection()
var db = config.GetDatabaseConnection()

func CurrentUser(c *gin.Context) *utils.AccessDetails {
	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
		return nil
	}
	return au
}
