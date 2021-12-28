package endpoints

import (
	"time"

	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListCredits(c *gin.Context) {
	var credits []objects.Credit
	var queryParams QueryParams
	order_by := []string{"created_at desc"}
	c.BindQuery(&queryParams)

	query := db.Connection.Preload("CreditType").Preload("User").Preload("Account").Model(&objects.Credit{}).Where("user_id = ? and account_id = ?", CurrentUser(c).UserId, queryParams.AccountId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&credits))
}

func CreateCredit(c *gin.Context) {
	var creditParams CreditParams
	c.ShouldBindJSON(&creditParams)

	credit := objects.Credit{
		CreditTypeID: creditParams.CreditTypeID,
		Amount:       creditParams.Amount,
		CreditedOn:   (*time.Time)(&creditParams.CreditedOn),
		UserID:       CurrentUser(c).UserId,
		AccountID:    creditParams.AccountID,
		Description:  creditParams.Description,
	}

	if err := db.Connection.Create(&credit).Error; err != nil {
		HandleError(c, err)
		return
	}
	db.Connection.Preload("CreditType").Preload("User").Preload("Account").Find(&credit)
	HandleSuccess(c, credit)
}
