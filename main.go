package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
)

func main() {

	fmt.Println("Bot starting...")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// DISCORD
	discordKey := os.Getenv("BOT_TOKEN")

	session, err := discordgo.New("Bot " + discordKey)
	if err != nil {
		fmt.Println("issue starting bot: ", err)
	}

	// OAI
	oaiKey := os.Getenv("OPEN_AI_KEY")
	model := openai.ChatModelGPT4_1Mini
	client := NewChatClient(oaiKey, model)

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handleMessage(s, m, client)
	})

	err = session.Open()
	if err != nil {
		panic(err)
	}

	select {}

}
