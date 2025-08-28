package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	Name        string `gorm:"index"`
	Designation string
	Email       string `gorm:"index"`
	Phone       string
	ImageURL    string
	ListID      uint    `gorm:"index"`
	CardOrder   float64 `gorm:"type:decimal(20,10);index"`
	Location    string

	List List  `gorm:"foreignKey:ListID;references:ID"`
	Tags []Tag `gorm:"many2many:card_tags;"`
}
