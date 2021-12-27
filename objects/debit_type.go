package objects

type DebitType struct {
	Base
	Name   string `gorm:"unique_index:idx_debit_types_name_user_id;index;not null;size:256" json:"name" binding:"required"`
	UserID string `gorm:"unique_index:idx_debit_types_name_user_id;index;not null" json:"user_id" binding:"required"`
	User   User
}
