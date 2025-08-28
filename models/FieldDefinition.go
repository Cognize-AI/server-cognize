package models

import "gorm.io/gorm"

type FieldDefinitionType string

const (
	CardTypeContact FieldDefinitionType = "CONTACT"
	CardTypeCompany FieldDefinitionType = "COMPANY"
)

func (t FieldDefinitionType) IsFieldTypeValid() bool {
	return t == CardTypeContact || t == CardTypeCompany
}

type FieldDefinition struct {
	gorm.Model
	Name     string
	DataType string `gorm:"default:'string'"`
	UserID   uint
	Type     string `gorm:"type:varchar(20)"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}
