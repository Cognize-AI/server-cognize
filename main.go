package main

import (
	"log"
	"time"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/internal/card"
	"github.com/Cognize-AI/client-cognize/internal/list"
	"github.com/Cognize-AI/client-cognize/internal/oauth"
	"github.com/Cognize-AI/client-cognize/internal/tag"
	"github.com/Cognize-AI/client-cognize/internal/user"
	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/router"
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
	config.SyncDB()
	logger.Logger.Info("DB sync completed")
}

func main() {
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {

		}
	}(logger.Logger)

	go func() {
		ticker := time.NewTicker(5 * time.Second) // flush every 5s
		defer ticker.Stop()
		for range ticker.C {
			if err := logger.Logger.Sync(); err != nil {
				logger.Logger.Error("failed to sync logs", zap.Error(err))
			} else {
				logger.Logger.Debug("log buffer flushed to Axiom")
			}
		}
	}()

	userSvc := user.NewService()
	oauthSvc := oauth.NewService()
	listSvc := list.NewService()
	cardSvc := card.NewService()
	tagSvc := tag.NewService()

	userHandler := user.NewHandler(userSvc)
	oauthHandler := oauth.NewHandler(oauthSvc)
	listHandler := list.NewHandler(listSvc)
	cardHandler := card.NewHandler(cardSvc)
	tagHandler := tag.NewHandler(tagSvc)

	router.InitRouter(
		userHandler,
		oauthHandler,
		listHandler,
		cardHandler,
		tagHandler,
	)
	log.Fatal(router.Start("0.0.0.0:" + Config.PORT))
}
