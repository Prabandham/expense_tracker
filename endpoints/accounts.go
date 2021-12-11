package endpoints

import (
	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListAccounts(c *gin.Context) {
	var accounts []objects.Account
	var queryParams QueryParams
	order_by := []string{"name asc"}
	c.BindQuery(&queryParams)

	query := db.Connection.Preload("User").Model(&objects.Account{}).Where("user_id = ?", CurrentUser(c).UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&accounts))
}

func CreateAccount(c *gin.Context) {
	var accountParams AccountParams
	c.BindJSON(&accountParams)

	account := objects.Account{
		Name:             accountParams.Name,
		Address:          accountParams.Address,
		IfscCode:         accountParams.IfscCode,
		AvailableBalance: accountParams.AvailableBalance,
		UserID:           CurrentUser(c).UserId,
	}

	if err := db.Connection.Create(&account).Error; err != nil {
		HandleError(c, err)
		return
	}
	HandleSuccess(c, &account)
}

func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	var account objects.Account
	if err := db.Connection.Where("id = ?", id).Delete(&account).Error; err != nil {
		HandleError(c, err)
		return
	}
	HandleSuccess(c, &account)
}
