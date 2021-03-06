package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Credit struct {
	Base
	CreditTypeID string `gorm:"index;not null" binding:"required" json:"credit_type_id"`
	CreditType   CreditType
	Description  string `gorm:"index;" json:"description"`
	UserID       string `gorm:"index;not null" binding:"required" json:"user_id"`
	User         User
	Amount       int64  `gorm:"not null;default:0" binding:"required;" json:"amount"`
	AccountID    string `gorm:"index;not null" binding:"required" json:"account_id"`
	Account      Account
	CreditedOn   *time.Time `sql:"index; not null" json:"credited_on"`
}

func (c *Credit) AfterSave(tx *gorm.DB) (err error) {
	var account Account
	tx.First(&account, "id = ?", c.AccountID)
	availableBalance := account.AvailableBalance
	newAvailableBalance := c.Amount + availableBalance
	account.AvailableBalance = newAvailableBalance
	tx.Save(&account)
	return nil
}