package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	Name        string `gorm:"index"`
	Designation string
	Email       string `gorm:"index"`
	Phone       string
	ImageURL    string
	ListID      uint `gorm:"index"`

	List List `gorm:"foreignKey:ListID;references:ID"`
}
