package endpoints

import (
	"time"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListDebits(c *gin.Context) {
	var debits []objects.Debit
	var queryParams QueryParams
	order_by := []string{"created_at desc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	db := config.GetDatabaseConnection()
	query := db.Connection.Preload("DebitType").Preload("User").Preload("Account").Model(&objects.Debit{}).Where("user_id = ?", au.UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&debits))
}

func CreateDebit(c *gin.Context) {
	var debitParams DebitParams
	HandleError(c, c.ShouldBindJSON(&debitParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	debit := objects.Debit{
		DebitTypeID: debitParams.DebitTypeID,
		Amount:   debitParams.Amount,
		DebitedOn:       (*time.Time)(&debitParams.DebitedOn),
		UserID:        au.UserId,
		AccountID: debitParams.AccountID,
	}

	db := config.GetDatabaseConnection()
	result := db.Connection.Create(&debit)
	db.Connection.Preload("DebitType").Preload("User").Preload("Account").Find(&debit)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, debit)
}
