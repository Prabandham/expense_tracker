package objects

import "time"

type Expense struct {
	Base
	ExpenseTypeID string
	ExpenseType ExpenseType
	UserID string
	AmountSpent int64 `gorm:"not null;default:0" binding:"required;" json:"amount_spent"`
	SpentOn *time.Time `sql:"index; not null" json:"spent_on"`
}