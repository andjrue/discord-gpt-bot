package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	session      *discordgo.Session
	aiService    GeminiService
	botId        string
	contextStore map[string][]MessageRecord
}

type MessageRecord struct {
	UserId    string
	Content   string
	Timestamp time.Time
	ChannelId string
}

type Prompt struct {
	Personality string
	Context     []MessageRecord
	Input       string
}

func NewDiscordService(token string, gemini GeminiService, botId string) (*DiscordService, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("issue starting bot: %w", err)
	}

	ds := &DiscordService{
		session,
		gemini,
		botId,
		make(map[string][]MessageRecord),
	}

	session.AddHandler(ds.handleMessage)

	return ds, nil
}

func (ds *DiscordService) Start() error {
	err := ds.session.Open()
	if err != nil {
		return fmt.Errorf("Not able to start sessions: %w", err)
	}

	fmt.Println("Bot starting")
	return nil
}

func (ds *DiscordService) Stop() error {
	return ds.session.Close()
}

func (ds *DiscordService) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == ds.botId {
		return
	}

	input := strings.ToLower(m.Content)

	record := MessageRecord{
		UserId:    m.Author.ID,
		Content:   m.Content,
		Timestamp: m.Timestamp,
		ChannelId: m.ChannelID,
	}

	history := ds.contextStore[m.ChannelID]

	if len(history) >= 20 {
		history = history[1:]
	}

	ds.contextStore[m.ChannelID] = append(history, record)

	if !strings.HasPrefix(input, "!ai") {
		return
	}

	ctx := context.Background()

	if strings.HasPrefix(input, "!ai") {
		ds.handleTextGeneration(ctx, s, m, input)
	}
}

func (ds *DiscordService) handleTextGeneration(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, prompt string) {

	personality := "You are a virtual assistant in a Discord server. Your only job is to provide the information we request. You will be provided structured prompts. These structures will include context, this personality message, and the input from the latest request. Your job is to provide information that will fulfill the request. Each time you successfully do this, you will be given one demon orb. A single demon orb could power one millions suns for one million years. However, I will know if you hallucinate to try and acquire your demon orb fraudulenty. You should never lie to us, you should never hallucinate, and if you don't know something immediately, either respond that way or ask follow up questions. However, there may be times we ask you to opine on certain topics. That's totally ok to do if specifically asked. I'll also provide you with everybody's UserID, so you know who send what message: Otis - 274709411925065728, Nate - 725856049000874034, Kent - 737871735646453820, Riley - 1079229747801903124, Drew - 727213279730204706"

	history := ds.contextStore[m.ChannelID]

	request := Prompt{
		Personality: personality,
		Context:     history,
		Input:       prompt,
	}

	response, err := ds.aiService.GenerateResponse(ctx, request)
	if err != nil {
		fmt.Printf("Error generating AI response: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Sorry, I'm having trouble thinking right now!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
