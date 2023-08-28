package main

import (
	"log"
	"os"

	"report-bot/service"

	//"github.com/joho/godotenv"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("failed to load .env: %s", err.Error())
	// }

	botToken := os.Getenv("BOT_TOKEN")

	h, err := service.New(botToken)
	if err != nil {
		log.Fatalf("failed to create botAPI: %s", err.Error())
	}

	if err := h.Run(); err != nil {
		log.Fatalf("failed to run bot: %s", err.Error())
	}
}
