package objects

type Account struct {
	Base
	Name             string `gorm:"unique_index:idx_accounts_name_user_id;index;not null;size:256" json:"name" binding:"required"`
	IfscCode         string `gorm:"unique_index:idx_accounts_name_user_id;index;not null;size:256" json:"ifsc_code"`
	Address          string `json:"address"`
	AvailableBalance int64  `gorm:"not null;default:0" json:"available_balance"`
	UserID           string `gorm:"unique_index:idx_accounts_name_user_id;index;not null;"`
	User             User
	Credits          []Credit
	Debits           []Debit
}

type AccountCreditsAndDebits struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	PerformedOn string `json:"performed_on"`
	Type        string `json:"type"`
	UserID      string `json:"user_id"`
	AccountID   string `json:"account_id"`
	Kind        string `json:"kind"`
}
