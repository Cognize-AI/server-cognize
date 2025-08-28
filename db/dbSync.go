package db

import (
	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/models"
)

func SyncDB() {
	err := config.DB.AutoMigrate(
		models.User{},
		models.List{},
		models.Card{},
		models.Tag{},
		models.Key{},
		models.Company{},
		models.Activity{},
		models.FieldDefinition{},
		models.FieldValue{},
	)
	if err != nil {
		return
	}
}
