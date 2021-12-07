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

type DebitTypeParams struct {
	Name string `json:"name" binding:"required"`
}
type CreditTypeParams struct {
	Name string `json:"name" binding:"required"`
}

type DebitParams struct {
	DebitTypeID string `json:"debit_type_id" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
	DebitedOn       myTime `json:"debited_on" binding:"required"`
}

type CreditParams struct {
	CreditTypeID string `json:"credit_type_id" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
	CreditedOn       myTime `json:"credited_on" binding:"required"`
}

type AccountParams struct {
	Name string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	IfscCode string `json:"ifsc_code" binding:"required"`
	AvailableBalance int64 `json:"available_balance"`
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
