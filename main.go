package main

import (
	"log"
	"time"

	"github.com/Cognize-AI/server-cognize/config"
	"github.com/Cognize-AI/server-cognize/db"
	"github.com/Cognize-AI/server-cognize/internal/card"
	"github.com/Cognize-AI/server-cognize/internal/keys"
	"github.com/Cognize-AI/server-cognize/internal/list"
	"github.com/Cognize-AI/server-cognize/internal/oauth"
	"github.com/Cognize-AI/server-cognize/internal/tag"
	"github.com/Cognize-AI/server-cognize/internal/user"
	"github.com/Cognize-AI/server-cognize/logger"
	"github.com/Cognize-AI/server-cognize/router"
	"go.uber.org/zap"
)

var Config config.Config

func init() {
	var err error
	Config, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(Config)
	logger.Logger.Info("Logger initialized")
	config.ConnectDB()
	logger.Logger.Info("DB connection established")
	db.SyncDB()
	logger.Logger.Info("DB sync completed")
}

func main() {
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {

		}
	}(logger.Logger)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = logger.Logger.Sync()
		}
	}()

	userSvc := user.NewService()
	oauthSvc := oauth.NewService()
	listSvc := list.NewService()
	cardSvc := card.NewService()
	tagSvc := tag.NewService()
	keySvc := keys.NewService()

	userHandler := user.NewHandler(userSvc)
	oauthHandler := oauth.NewHandler(oauthSvc)
	listHandler := list.NewHandler(listSvc)
	cardHandler := card.NewHandler(cardSvc)
	tagHandler := tag.NewHandler(tagSvc)
	keyHandler := keys.NewHandler(keySvc)

	router.InitRouter(
		userHandler,
		oauthHandler,
		listHandler,
		cardHandler,
		tagHandler,
		keyHandler,
	)
	log.Fatal(router.Start("0.0.0.0:" + Config.PORT))
}
