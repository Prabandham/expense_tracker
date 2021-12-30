package endpoints

import (
	"github.com/Prabandham/expense_tracker/utils"
	"github.com/gin-gonic/gin"
)

func GroupedCredits(c *gin.Context) {
	type Result struct {
		Amount int
		Type   string
	}
	result := []Result{}
	sql := QueryRepo().Get("groupedCredits")
	db.Connection.Raw(sql, CurrentUser(c).UserId).Scan(&result)
	enc := utils.NewStructEncoder()
	rows, err := enc.Marshal(result)
	if err != nil {
		panic(err)
	}
	HandleSuccess(c, &rows)
}

func GroupedDebits(c *gin.Context) {
	type Result struct {
		Amount int
		Type   string
	}
	result := []Result{}
	sql := QueryRepo().Get("groupedDebits")
	db.Connection.Raw(sql, CurrentUser(c).UserId).Scan(&result)
	enc := utils.NewStructEncoder()
	rows, err := enc.Marshal(result)
	if err != nil {
		panic(err)
	}
	HandleSuccess(c, &rows)
}
