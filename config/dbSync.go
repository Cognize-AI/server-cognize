package config

import "github.com/Cognize-AI/client-cognize/models"

func SyncDB() {
	err := DB.AutoMigrate(
		models.User{},
		models.List{},
		models.Card{},
		models.Tag{},
	)
	if err != nil {
		return
	}
}
