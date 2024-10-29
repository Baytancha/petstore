package main

import (
	"log"
	"os"
	"test/config"
	"test/run"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found here")
	}
}

func main() {
	cfg := config.NewConfig(config.WithPort(8080), config.WithDBname("postgres"), config.WithDSN(os.Getenv("DB_DSN")))

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := run.NewApp(cfg, logger)

	app.Run()

	err = app.Serve()
	logger.Error(err.Error())
	os.Exit(1)

}
