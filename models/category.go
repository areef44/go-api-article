package models

type Category struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100)" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}
