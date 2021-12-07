package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListCreditTypes(c *gin.Context) {
	var CreditTypes []objects.CreditType
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
	HandleSuccess(c, paginator.Paginate(&CreditTypes))
}

func CreateCreditType(c *gin.Context) {
	var creditTypeParams CreditTypeParams
	HandleError(c, c.ShouldBindJSON(&creditTypeParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	db := config.GetDatabaseConnection()
	credit_type := objects.CreditType{Name: creditTypeParams.Name, UserID: au.UserId}
	result := db.Connection.Create(&credit_type)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, result.Value)
}
