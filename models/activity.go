package models

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	Content string
	CardID  uint `gorm:"not null;uniqueIndex"`

	Card Card `gorm:"foreignKey:CardID;references:ID"`
}
