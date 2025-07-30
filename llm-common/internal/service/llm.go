package service

import (
	"context"
	"github.com/compendium-tech/compendium/llm-common/internal/domain"
)

type LLMService interface {
	GenerateResponse(
		ctx context.Context, chatHistory []domain.Message,
		structuredOutputSchema *domain.StructuredOutputSchema) (*domain.Message, error)
}
