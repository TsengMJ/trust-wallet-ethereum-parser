package main

import (
	"ethereum-parser/config"
	"ethereum-parser/cron"
	"ethereum-parser/logger"
	pubsub "ethereum-parser/pkg/pub-sub"
	"ethereum-parser/server"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Load config
	var configPath = os.Getenv("CONFIG_PATH")
	var env = os.Getenv("ENV")

	err = config.InitConfig(&configPath, &env)
	if err != nil {
		panic("Error loading config file, " + err.Error())
	}

	err = logger.InitLog()
	if err != nil {
		panic("Error initializing logger, " + err.Error())
	}

	publisher := pubsub.DefaultPublisher
	cron := cron.NewListenEthereumBlockCron(publisher)
	go cron.Start()

	// Start rest api server
	server.StartServer()
}
