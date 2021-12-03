package endpoints

import (
	"time"
	"encoding/json"
)

type LoginParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegesterParams struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ExpenseTypeParams struct {
	Name string `json:"name" binding:"required"`
}

type ExpenseParams struct {
	ExpenseTypeID string `json:"expense_type_id" binding:"required"`
	AmountSpent   int64  `json:"amount_spent" binding:"required"`
	SpentOn       myTime `json:"spent_on" binding:"required"`
}

type QueryParams struct {
	Page    string `form:"page"`
	PerPage string `form:"per_page"`
}

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
