package objects

type CreditType struct {
	Base
	Name   string `gorm:"unique_index:idx_credit_types_name_user_id;index;not null;size:256" json:"name" binding:"required"`
	UserID string `gorm:"unique_index:idx_credit_types_name_user_id;index;not null" json:"user_id" binding:"required"`
	User   User
}
