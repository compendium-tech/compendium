package service

import (
	"context"
	"fmt"
	"time"

	"github.com/compendium-tech/compendium/llm-service/internal/domain"
	"google.golang.org/genai"
)

type geminiClient struct {
	client *genai.Client
	tools  []*genai.Tool
}

func NewGeminiClient(ctx context.Context, cfg *genai.ClientConfig) (LLMService, error) {
	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &geminiClient{client: client}, nil
}

func (g *geminiClient) GenerateResponse(
	ctx context.Context,
	chatHistory []domain.Message,
	tools []domain.ToolDefinition,
	structuredOutputSchema *domain.Schema,
) (*domain.Message, error) {
	g.useTools(tools)

	contents := make([]*genai.Content, 0, len(chatHistory))
	for _, msg := range chatHistory {
		parts := make([]*genai.Part, 0, 1+len(msg.ToolCalls))
		parts = append(parts, &genai.Part{Text: msg.Text})

		for _, toolCall := range msg.ToolCalls {
			parts = append(parts, &genai.Part{
				FunctionCall: &genai.FunctionCall{
					Name: toolCall.Name,
					Args: toolCall.Parameters,
				},
			})
		}
		contents = append(contents, &genai.Content{
			Parts: parts,
			Role:  string(msg.Role),
		})
	}

	var schema *genai.Schema
	if structuredOutputSchema != nil {
		schema = domainSchemaToGenAISchema(structuredOutputSchema)
	}

	config := &genai.GenerateContentConfig{
		ResponseSchema: schema,
		Tools:          g.tools,
		Temperature:    genai.Ptr[float32](0),
	}

	result, err := g.client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	message := domain.Message{
		Role: domain.RoleAssistant,
	}

	text := result.Text()
	message.Text = text

	for _, candidate := range result.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.FunctionCall != nil {
				toolCall := domain.ToolCall{
					ID:         fmt.Sprintf("%v%d", part.FunctionCall.Name, time.Now().UnixNano()),
					Name:       part.FunctionCall.Name,
					Parameters: part.FunctionCall.Args,
				}
				message.ToolCalls = append(message.ToolCalls, toolCall)
			}
		}
	}

	return &message, nil
}

func domainSchemaToGenAISchema(domainSchema *domain.Schema) *genai.Schema {
	if domainSchema == nil {
		return nil
	}

	ty := toGenAIType(domainSchema.Type)
	properties := make(map[string]*genai.Schema)
	for name, prop := range domainSchema.Properties {
		properties[name] = domainSchemaToGenAISchema(&prop)
	}

	var items *genai.Schema
	if domainSchema.Items != nil {
		items = domainSchemaToGenAISchema(domainSchema.Items)
	}

	return &genai.Schema{
		Type:        ty,
		Description: domainSchema.Description,
		Properties:  properties,
		Items:       items,
		MaxItems:    domainSchema.MaxItems,
		MinItems:    domainSchema.MinItems,
		Required:    domainSchema.Required,
	}
}

func toGenAIType(t domain.Type) genai.Type {
	switch t {
	case domain.TypeString:
		return genai.TypeString
	case domain.TypeNumber:
		return genai.TypeNumber
	case domain.TypeInteger:
		return genai.TypeInteger
	case domain.TypeBoolean:
		return genai.TypeBoolean
	case domain.TypeArray:
		return genai.TypeArray
	case domain.TypeObject:
		return genai.TypeObject
	default:
		return genai.TypeString
	}
}

func (g *geminiClient) useTools(tools []domain.ToolDefinition) {
	genAITools := make([]*genai.Tool, 0, len(tools))
	for _, toolDef := range tools {
		genAITools = append(genAITools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        toolDef.Name,
					Description: toolDef.Description,
					Parameters:  domainSchemaToGenAISchema(toolDef.ParametersSchema),
				},
			},
		})
	}

	g.tools = genAITools
}
