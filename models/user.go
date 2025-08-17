package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"uniqueIndex"`
	Password       string
	ProfilePicture string

	Lists []List `gorm:"foreignKey:UserID;references:ID"`
}
