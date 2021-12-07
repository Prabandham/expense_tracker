package objects

type Account struct {
	Base
	Name string `gorm:"index;not null;size:256" json:"name" binding:"required"`
	IfscCode string `gorm:"index;not null;size:256" json:"ifsc_code"`
	Address string `json:"address"`
	AvailableBalance int64 `gorm:"not null;default:0" json:"available_balance"`
	UserID string `gorm:"index;not null;"`
	User User
	Credits []Credit
	Debits []Debit
}