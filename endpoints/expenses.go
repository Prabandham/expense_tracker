package endpoints

import (
	"time"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListExpenses(c *gin.Context) {
	var expenses []objects.Expense
	var queryParams QueryParams
	order_by := []string{"created_at desc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	db := config.GetDatabaseConnection()
	query := db.Connection.Preload("ExpenseType").Preload("User").Model(&objects.ExpenseType{}).Where("user_id = ?", au.UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&expenses))
}

func CreateExpense(c *gin.Context) {
	var expenseParams ExpenseParams
	HandleError(c, c.ShouldBindJSON(&expenseParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	expense := objects.Expense{
		ExpenseTypeID: expenseParams.ExpenseTypeID,
		AmountSpent:   expenseParams.AmountSpent,
		SpentOn:       (*time.Time)(&expenseParams.SpentOn),
		UserID:        au.UserId,
	}

	db := config.GetDatabaseConnection()
	result := db.Connection.Create(&expense)
	db.Connection.Preload("ExpenseType").Preload("User").Find(&expense)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, expense)
}
