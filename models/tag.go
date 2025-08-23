package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name   string
	Color  string
	UserID uint `gorm:"index"`

	User  User   `gorm:"foreignKey:UserID;references:ID"`
	Cards []Card `gorm:"many2many:card_tags;"`
}
