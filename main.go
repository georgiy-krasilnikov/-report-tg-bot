package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"report-bot/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	botToken := os.Getenv("BOT_TOKEN")

	h, err := services.New(botToken)
	if err != nil {
		panic(err)
	}

	if err := h.Run(); err != nil {
		log.Fatalf("failed to run bot: %s", err.Error())
	}
}
