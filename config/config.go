package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordKey string
	GeminiKey  string
	Model      string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	discordToken := os.Getenv("BOT_TOKEN")
	if discordToken == "" {
		return nil, fmt.Errorf("Unable to load Discord token")
	}

	geminiKey := os.Getenv("GEMINI_KEY")
	if geminiKey == "" {
		return nil, fmt.Errorf("Unable to load Gemini key")
	}

	model := "gemini-2.5-pro"

	return &Config{
		discordToken,
		geminiKey,
		model,
	}, nil
}
