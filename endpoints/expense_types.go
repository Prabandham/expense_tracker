package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListExpenseTypes(c *gin.Context) {
	var ExpenseTypes []objects.ExpenseType
	var queryParams QueryParams
	order_by := []string{"name asc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	db := config.GetDatabaseConnection()
	paginator := p.Paginator{DB: db.Connection, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&ExpenseTypes))
}

func CreateExpenseType(c *gin.Context) {
	var expenseTypeParams ExpenseTypeParams
	HandleError(c, c.ShouldBindJSON(&expenseTypeParams))

	db := config.GetDatabaseConnection()
	expense_type := objects.ExpenseType{Name: expenseTypeParams.Name}
	result := db.Connection.Create(&expense_type)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, result.Value)
}
