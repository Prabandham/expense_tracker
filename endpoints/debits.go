package endpoints

import (
	"time"

	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListDebits(c *gin.Context) {
	var queryParams QueryParams
	c.ShouldBindQuery(&queryParams)
	order_by := []string{"created_at desc"}
	var debits []objects.Debit

	query := db.Connection.Preload("DebitType").Preload("User").Preload("Account").Model(&objects.Debit{}).Where("user_id = ?", CurrentUser(c).UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&debits))
}

func CreateDebit(c *gin.Context) {
	var debitParams DebitParams
	c.ShouldBindJSON(&debitParams)

	debit := objects.Debit{
		DebitTypeID: debitParams.DebitTypeID,
		Amount:      debitParams.Amount,
		DebitedOn:   (*time.Time)(&debitParams.DebitedOn),
		UserID:      CurrentUser(c).UserId,
		AccountID:   debitParams.AccountID,
	}

	if err := db.Connection.Create(&debit).Error; err != nil {
		HandleError(c, err)
		return
	}
	db.Connection.Preload("DebitType").Preload("User").Preload("Account").Find(&debit)
	HandleSuccess(c, debit)
}
