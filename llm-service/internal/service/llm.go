package service

import (
	"context"

	"github.com/compendium-tech/compendium/llm-service/internal/domain"
)

type LLMService interface {
	GenerateResponse(
		ctx context.Context, chatHistory []domain.Message,
		tools []domain.ToolDefinition, structuredOutputSchema *domain.Schema) (*domain.Message, error)
}
