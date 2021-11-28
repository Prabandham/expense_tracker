package endpoints

import (
	"net/http"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListExpenseTypes(c *gin.Context) {
	var ExpenseTypes []objects.ExpenseType
	var queryParams QueryParams
	order_by := []string{"name asc"}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	db := config.GetDatabaseConnection()
	paginator := p.Paginator{DB: db.Connection, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	data := paginator.Paginate(&ExpenseTypes)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func CreateExpenseType(c *gin.Context) {
	var expenseParams ExpenseTypeParams
	db := config.GetDatabaseConnection()
	if err := c.ShouldBindJSON(&expenseParams); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	expense_type := objects.ExpenseType{Name: expenseParams.Name}
	result := db.Connection.Create(&expense_type)
	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": expense_type,
	})
}

type ExpenseTypeParams struct {
	Name string `json:"name" binding:"required"`
}

type QueryParams struct {
	Page string `form:"page"`
	PerPage string `form:"per_page"`
}
