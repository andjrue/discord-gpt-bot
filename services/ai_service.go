package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService interface {
	GenerateResponse(ctx context.Context, prompt string) (string, error)
	SearchWeb(ctx context.Context, query string) (string, error)
	Close() error
}

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewService(apiKey string, modelName string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Unable to create client: %w", err)
	}

	model := client.GenerativeModel(modelName)

	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiService) GenerateResponse(ctx context.Context, p Prompt) (string, error) {

	var sb strings.Builder

	fmt.Println("Prompt: \n", p)

	sb.WriteString(p.Personality)
	sb.WriteString("\n\nContext:\n\n")

	for _, msg := range p.Context {
		sb.WriteString(fmt.Sprintf("<%s>: %s\n", msg.UserId, msg.Content))
	}

	sb.WriteString("\nUser Prompt: " + p.Input)

	finalPrompt := sb.String()

	fmt.Printf("Final Prompt: %s", finalPrompt)

	resp, err := g.model.GenerateContent(ctx, genai.Text(finalPrompt))
	if err != nil {
		return "", fmt.Errorf("Not able to generate response: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("No response generated")
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			result += string(txt)
		}
	}

	return result, nil
}

func (g *GeminiService) Close() error {
	return g.client.Close()
}
