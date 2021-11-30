package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Prabandham/expense_tracker/config"
	"github.com/Prabandham/expense_tracker/objects"
	"github.com/Prabandham/expense_tracker/utils"
	p "github.com/Prabandham/paginator"
	"github.com/gin-gonic/gin"
)

type myTime time.Time

var _ json.Unmarshaler = &myTime{}

func (mt *myTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02T15:04:05Z", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = myTime(t)
	return nil
}

func ListExpenses(c *gin.Context) {
	var expenses []objects.Expense
	var queryParams QueryParams
	order_by := []string{"created_at desc"}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	db := config.GetDatabaseConnection()
	au, err := utils.ExtractTokenMetadata(c.Request)
	db.Connection.Model(&objects.User{}).Where("id = ?", au.UserId).Find(&expenses)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	paginator := p.Paginator{DB: db.Connection, OrderBy: order_by, Page: queryParams.Page, PerPage: queryParams.PerPage}
	data := paginator.Paginate(&expenses)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func CreateExpense(c *gin.Context) {
	var params ExpenseParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	au, err := utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := config.GetDatabaseConnection()
	result := db.Connection.Create(&objects.Expense{
		ExpenseTypeID: params.ExpenseTypeID,
		AmountSpent:   params.AmountSpent,
		SpentOn:       (*time.Time)(&params.SpentOn),
		UserID:        au.UserId,
	})

	if result.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, result.Error)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

type ExpenseParams struct {
	ExpenseTypeID string `json:"expense_type_id" binding:"required"`
	AmountSpent   int64  `json:"amount_spent" binding:"required"`
	SpentOn       myTime `json:"spent_on" binding:"required"`
}
