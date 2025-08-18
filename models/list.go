package models

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Name      string
	Color     string
	UserID    uint    `gorm:"index"`
	ListOrder float64 `gorm:"type:decimal(10,2);index"`

	User  User   `gorm:"foreignKey:UserID;references:ID"`
	Cards []Card `gorm:"foreignKey:ListID;references:ID"`
}
