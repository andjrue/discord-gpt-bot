package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andjrue/discord-gpt-bot/config"
	"github.com/andjrue/discord-gpt-bot/services"
)

func main() {

	fmt.Println("Bot starting")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Errorf("Not able to load config: %w", err)
	}

	gemini, err := services.NewService(cfg.GeminiKey, cfg.Model)
	if err != nil {
		fmt.Errorf("Not able to load ai service: %w", err)
	}
	defer gemini.Close()

	botId := "1388713265110188182"
	discordService, err := services.NewDiscordService(cfg.DiscordKey, *gemini, botId)
	if err != nil {
		fmt.Errorf("Unable to load discord service: %w", err)
	}

	if err := discordService.Start(); err != nil {
		fmt.Errorf("Unable to start discord service: %w", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down bot...")

}
