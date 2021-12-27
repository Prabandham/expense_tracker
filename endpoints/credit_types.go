package endpoints

import (
	"github.com/Prabandham/expense_tracker/objects"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

func ListCreditTypes(c *gin.Context) {
	var CreditTypes []objects.CreditType
	var queryParams QueryParams
	order_by := []string{"name asc"}
	c.BindQuery(&queryParams)

	query := db.Connection.Preload("User").Where("user_id = ?", CurrentUser(c).UserId)
	paginator := p.Paginator{DB: query, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	HandleSuccess(c, paginator.Paginate(&CreditTypes))
}

func CreateCreditType(c *gin.Context) {
	var creditTypeParams CreditTypeParams
	c.ShouldBindJSON(&creditTypeParams)

	credit_type := objects.CreditType{Name: creditTypeParams.Name, UserID: CurrentUser(c).UserId}
	if err := db.Connection.Create(&credit_type).Error; err != nil {
		HandleError(c, err)
		return
	}

	HandleSuccess(c, &credit_type)
}

func DeleteCreditType(c *gin.Context) {
	id := c.Param("id")
	var creditType objects.CreditType
	if err := db.Connection.Where("id = ?", id).Delete(&creditType).Error; err != nil {
		HandleError(c, err)
		return
	}
	HandleSuccess(c, &creditType)
}
