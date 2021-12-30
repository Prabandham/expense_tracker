package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
	"github.com/Davmuz/gqt"
)

var redis = config.GetRedisConnection()
var db = config.GetDatabaseConnection()

// CurrentUser is a convenience function to extract the user
// from the auth token
func CurrentUser(c *gin.Context) *utils.AccessDetails {
	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
		return nil
	}
	return au
}

func QueryRepo() *gqt.Repository {
	sql := gqt.NewRepository()
	queryRepoPath := config.GetEnv("QUERY_REPO_PATH", "/Users/prabandham/projects/go/src/github.com/Prabandham/expense_tracker/config")
	sql.Add(queryRepoPath, "*.sql")
	return sql
}
