package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListDebitTypes(c *gin.Context) {
	var DebitTypes []objects.DebitType
	var queryParams QueryParams
	order_by := []string{"name asc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	db := config.GetDatabaseConnection()
	query := db.Connection.Preload("User").Where("user_id = ?", au.UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&DebitTypes))
}

func CreateDebitType(c *gin.Context) {
	var debitTypeParams DebitTypeParams
	HandleError(c, c.ShouldBindJSON(&debitTypeParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	db := config.GetDatabaseConnection()
	debit_type := objects.DebitType{Name: debitTypeParams.Name, UserID: au.UserId}
	result := db.Connection.Create(&debit_type)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, result.Value)
}
