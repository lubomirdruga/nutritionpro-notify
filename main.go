package main

import (
	"log"
	"nutritionpro-notify/telegram"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("TELEGRAM_API_TOKEN")
	if token == "" {
		panic("TELEGRAM_API_TOKEN environment variable is not set")
	}

	botService, err := telegram.NewBotService(token)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := botService.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for a shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	botService.Stop()
}
