package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isValidMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Author.ID == s.State.User.ID {
		return false
	}

	return strings.HasPrefix(m.Content, "!ai")
}

func extractPrompt(m *discordgo.MessageCreate) string {
	return strings.TrimPrefix(m.Content, "!ai")
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, client *ChatClient) {

	fmt.Println("Message Received: ", m.Content)

	if !isValidMessage(s, m) {
		return
	}

	userPrompt := extractPrompt(m)

	response, err := client.SendRequestToOai(userPrompt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong with the request.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
