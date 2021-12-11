package endpoints

import (
	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListDebitTypes(c *gin.Context) {
	var DebitTypes []objects.DebitType
	var queryParams QueryParams
	order_by := []string{"name asc"}
	c.ShouldBindQuery(&queryParams)

	query := db.Connection.Preload("User").Where("user_id = ?", CurrentUser(c).UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&DebitTypes))
}

func CreateDebitType(c *gin.Context) {
	var debitTypeParams DebitTypeParams
	c.ShouldBindJSON(&debitTypeParams)

	debit_type := objects.DebitType{Name: debitTypeParams.Name, UserID: CurrentUser(c).UserId}
	if err := db.Connection.Create(&debit_type).Error; err != nil {
		HandleError(c, err)
		return
	}
	HandleSuccess(c, &debit_type)
}

func DeleteDebitType(c *gin.Context) {
	id := c.Param("id")
	var debitType objects.DebitType
	if err := db.Connection.Where("id = ?", id).Delete(&debitType).Error; err != nil {
		HandleError(c, err)
		return
	}
	HandleSuccess(c, &debitType)
}