package config

import "github.com/Cognize-AI/client-cognize/models"

func SyncDB() {
	err := DB.AutoMigrate(
		models.User{},
		models.List{},
		models.Card{},
	)
	if err != nil {
		return
	}
}
