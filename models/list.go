package models

import "gorm.io/gorm"

type List struct {
	gorm.Model
	Name      string
	Color     string
	UserID    uint `gorm:"index"`
	ListOrder uint `gorm:"autoIncrement"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}
