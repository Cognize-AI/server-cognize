package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name     string
	Role     string
	Location string
	Phone    string
	Email    string
	CardID   uint `gorm:"not null;uniqueIndex"`

	Card Card `gorm:"foreignKey:CardID;references:ID"`
}
