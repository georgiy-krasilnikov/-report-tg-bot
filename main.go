package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"report-bot/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %s", err.Error())
	}

	botToken := os.Getenv("BOT_TOKEN")

	h, err := services.New(botToken)
	if err != nil {
		log.Fatalf("failed to create botAPI: %s", err.Error())
	}

	if err := h.Run(); err != nil {
		log.Fatalf("failed to run bot: %s", err.Error())
	}
}
