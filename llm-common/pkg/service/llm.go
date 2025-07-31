package service

import (
	"context"
	"github.com/compendium-tech/compendium/llm-common/pkg/domain"
)

type LLMService interface {
	GenerateResponse(
		ctx context.Context, chatHistory []domain.Message,
		tools []domain.ToolDefinition,
		structuredOutputSchema *domain.StructuredOutputSchema) (*domain.Message, error)
}
