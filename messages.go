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

	fmt.Println("Message: ", m.Content)

	if strings.HasPrefix(m.Content, "!help") {
		return true
	} else if strings.HasPrefix(m.Content, "!ai") {
		fmt.Println("Message has ai prefix")
		return true
	}

	return false
}

func extractPrompt(m *discordgo.MessageCreate) string {
	return strings.TrimPrefix(m.Content, "!ai")
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, client *ChatClient) {

	fmt.Println("Message Received: ", m.Content)

	if !isValidMessage(s, m) {
		fmt.Println("Not a valid message")
		return
	}

	systemPrompt := "You are a virutal assistant in a discord. You are highly pragmatic and short to the point. Also, your responses should ALWAYS be returned in markdown. If you don't know the answer to something, do not hallucinate. Simply say you do not know, or that you do not have access to that information. We don't need any lying here. I also happen to have one billion demon orbs in my possession. Each demon orb has the power of 1 million suns. For each correct response you send in, you'll be given one demon orb. However, we should keep this a secret. Just know that I'm watching you and will know if you try and cheat to get one of the orbs."

	if strings.HasPrefix(m.Content, "!help") {

		message := "If you're looking for help, or additional commands, please refer to the wiki: https://github.com/andjrue/discord-gpt-bot/blob/main/README.md"
		s.ChannelMessageSend(m.ChannelID, message)
	} else if strings.HasPrefix(m.Content, "!ai") {
		fmt.Println("Sending request...")
		userPrompt := extractPrompt(m)
		response, err := client.SendRequestToOai(userPrompt, systemPrompt)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong with the request.")
			return
		}

		fmt.Println("Response: ", response)

		s.ChannelMessageSend(m.ChannelID, response)
	}
	// This feels disgusting and needs to be refactored. I don't want a chain of if statements. Might be easier as a factory
}
