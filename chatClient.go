package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type ChatClient struct {
	client openai.Client
	model  shared.ChatModel
}

func NewChatClient(apiKey, model shared.ChatModel) *ChatClient {

	c := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ChatClient{
		client: c,
		model:  model,
	}
}

// update model here, not sure what the type is
func (client *ChatClient) SendRequestToOai(userMessage string) (string, error) {

	ctx := context.Background()

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userMessage),
		},
		Seed:  openai.Int(0),
		Model: client.model,
	}

	completion, err := client.client.Chat.Completions.New(ctx, params)

	if err != nil {
		panic(err)
	}

	res := completion.Choices[0].Message.Content

	return res, nil
}
