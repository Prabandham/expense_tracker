package objects

type ExpenseType struct {
	Base
	Name string `gorm:"unique;index;not null;size:256" json:"name" binding:"required"`
}