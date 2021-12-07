package endpoints

import (
	"time"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListCredits(c *gin.Context) {
	var credits []objects.Credit
	var queryParams QueryParams
	order_by := []string{"created_at desc"}
	HandleError(c, c.ShouldBindQuery(&queryParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	HandleError(c, err)

	db := config.GetDatabaseConnection()
	query := db.Connection.Preload("CreditType").Preload("User").Preload("Account").Model(&objects.Credit{}).Where("user_id = ?", au.UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&credits))
}

func CreateCredit(c *gin.Context) {
	var creditParams CreditParams
	HandleError(c, c.ShouldBindJSON(&creditParams))

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		HandleError(c, err)
	}

	credit := objects.Credit{
		CreditTypeID: creditParams.CreditTypeID,
		Amount:   creditParams.Amount,
		CreditedOn:       (*time.Time)(&creditParams.CreditedOn),
		UserID:        au.UserId,
		AccountID: creditParams.AccountID,
	}

	db := config.GetDatabaseConnection()
	result := db.Connection.Create(&credit)
	db.Connection.Preload("CreditType").Preload("User").Preload("Account").Find(&credit)
	if result.Error != nil {
		HandleError(c, result.Error)
	}

	HandleSuccess(c, credit)
}
