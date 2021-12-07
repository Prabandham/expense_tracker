package endpoints

import (
	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListAccounts(c *gin.Context) {
	var accounts []objects.Account
	var queryParams QueryParams
	order_by := []string{"name asc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	db := config.GetDatabaseConnection()
	query := db.Connection.Preload("User").Model(&objects.Account{}).Where("user_id = ?", au.UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&accounts))
}

func CreateAccount(c *gin.Context) {
	var accountParams AccountParams
	HandleError(c, c.ShouldBindJSON(&accountParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	account := objects.Account{
		Name:   accountParams.Name,
		Address: accountParams.Address,
		IfscCode: accountParams.IfscCode,
		AvailableBalance: accountParams.AvailableBalance,
		UserID: au.UserId,
	}

	db := config.GetDatabaseConnection()
	result := db.Connection.Create(&account)
	db.Connection.Preload("User").Find(&account)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, account)
}