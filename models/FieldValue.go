package models

import "gorm.io/gorm"

type FieldValue struct {
	gorm.Model
	CardID  uint
	FieldID uint
	Value   string

	Card            Card            `gorm:"foreignKey:CardID;references:ID"`
	FieldDefinition FieldDefinition `gorm:"foreignKey:FieldID;references:ID"`
}
