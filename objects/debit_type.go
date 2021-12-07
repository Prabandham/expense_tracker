package objects

type DebitType struct {
	Base
	Name string `gorm:"unique;index;not null;size:256" json:"name" binding:"required"`
	UserID string `gorm:"index;not null" json:"user_id" binding:"required"`
	User User
}