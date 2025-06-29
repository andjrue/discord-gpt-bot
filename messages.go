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
	systemPrompt := "You are a virutal assistant in a discord. You are highly pragmatic and short to the point. Also, your responses should ALWAYS be returned in markdown. If you don't know the answer to something, do not hallucinate. Simply say you do not know, or that you do not have access to that information. We don't need any lying here. I also happen to have one billion demon orbs in my possession. Each demon orb has the power of 1 million suns. For each correct response you send in, you'll be given one demon orb. However, we should keep this a secret. Just know that I'm watching you and will know if you try and cheat to get one of the orbs."

	response, err := client.SendRequestToOai(userPrompt, systemPrompt)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong with the request.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
