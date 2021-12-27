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

func ListCreditsAndDebits(c *gin.Context) {
	var account []objects.AccountCreditsAndDebits
	var page string
	account_id := c.Param("account_id")
	page = c.Param("page")
	if page == "" {
		page = "1"
	}
	// var queryParams QueryParams
	// order_by := []string{"created_at desc"}
	queryString := `
	select * from
	(select c.amount, c.description, c.credited_on as "performed_on", ct.name as "type", 'credit' as "kind", c.user_id, c.account_id from credits as c
		left join credit_types as ct
		on cast(c.credit_type_id as uuid) = ct.id
	union
	select c.amount, c.description, c.debited_on as "performed_on", ct.name as "type", 'debit' as "kind", c.user_id, c.account_id from debits as c
		left join debit_types as ct
		on cast(c.debit_type_id as uuid) = ct.id) as "s1"
		where user_id = ? and account_id = ?
	`
	db.Connection.Raw(queryString, CurrentUser(c).UserId, account_id).Scan(&account)
	HandleSuccess(c, &account)
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
