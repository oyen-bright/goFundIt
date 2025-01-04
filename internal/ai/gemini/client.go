package gemini

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/oyen-bright/goFundIt/internal/ai/interfaces"
	"google.golang.org/api/option"
)

type geminiClient struct {
	client  *genai.Client
	context *context.Context
	model   *genai.GenerativeModel
}

func NewClient(apiKey string) (interfaces.AIService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	return &geminiClient{
		context: &ctx,
		client:  client,
		model:   model,
	}, nil
}

func Close(service interfaces.AIService) {
	service.(*geminiClient).client.Close()
}
