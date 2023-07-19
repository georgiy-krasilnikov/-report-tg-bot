package main

import (
	"os"

	"github.com/joho/godotenv"

	"report-bot/services"
)

func main() {
	if err := godotenv.Load(); err != nil {

	}

	botToken := os.Getenv("BOT_TOKEN")

	h, err := services.New(botToken)
	if err != nil {
		panic(err)
	}

	if err := h.Run(); err != nil {
		panic(err)
	}
}
