package objects

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Debit struct {
	Base
	DebitTypeID string `gorm:"index;not null" binding:"required" json:"debit_type_id"`
	DebitType   DebitType
	Description string `gorm:"index;" json:"description"`
	UserID      string
	User        User
	Amount      int64  `gorm:"not null;default:0" binding:"required;" json:"amount"`
	AccountID   string `gorm:"index;not null" binding:"required" json:"account_id"`
	Account     Account
	DebitedOn   *time.Time `sql:"index; not null" json:"debited_on"`
}

func (d *Debit) AfterSave(tx *gorm.DB) (err error) {
	var account Account
	tx.First(&account, "id = ?", d.AccountID)
	availableBalance := account.AvailableBalance
	newAvailableBalance := availableBalance - d.Amount
	account.AvailableBalance = newAvailableBalance
	tx.Save(&account)
	return nil
}