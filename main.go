package main

import (
	"fmt"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/Cognize-AI/client-cognize/logger"
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
	defer logger.Logger.Sync()
}

func main() {
	fmt.Println("Hello World")
}
