package service

import (
	"context"
	"fmt"
	"time"

	"github.com/compendium-tech/compendium/llm-common/pkg/domain"
	"google.golang.org/genai"
)

type geminiClient struct {
	client *genai.Client
	tools  []*genai.Tool
}

func NewGeminiClient(ctx context.Context, cfg *genai.ClientConfig, tools []domain.ToolDefinition) (LLMService, error) {
	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	c := &geminiClient{client: client}
	c.addTools(tools)

	return c, nil
}

func (g *geminiClient) addTools(tools []domain.ToolDefinition) {
	genAITools := make([]*genai.Tool, 0, len(tools))
	for _, toolDef := range tools {
		parameters := make([]*genai.Schema, len(toolDef.Parameters))
		for i, param := range toolDef.Parameters {
			parameters[i] = &genai.Schema{
				Type:        toGenAIType(param.Type),
				Description: param.Description,
				Enum:        param.Enum,
			}
		}

		parameterSchema := map[string]any{
			"type":       "object",
			"properties": make(map[string]any),
			"required":   make([]string, 0),
		}
		properties := parameterSchema["properties"].(map[string]any)
		for _, param := range toolDef.Parameters {
			m := map[string]any{
				"type":        string(param.Type),
				"description": param.Description,
			}

			if len(param.Enum) > 0 {
				m["enum"] = param.Enum
			}

			properties[param.Name] = m

			if param.IsRequired {
				parameterSchema["required"] = append(parameterSchema["required"].([]string), param.Name)
			}
		}

		genAITools = append(genAITools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:                 toolDef.Name,
					Description:          toolDef.Description,
					ParametersJsonSchema: parameterSchema,
				},
			},
		})
	}

	g.tools = genAITools
}

func (g *geminiClient) GenerateResponse(
	ctx context.Context,
	chatHistory []domain.Message,
	structuredOutputSchema *domain.StructuredOutputSchema,
) (*domain.Message, error) {
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
		schema = structuredOutputSchemaToGenAI(structuredOutputSchema)
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

func structuredOutputSchemaToGenAI(domainSchema *domain.StructuredOutputSchema) *genai.Schema {
	if domainSchema == nil {
		return nil
	}

	ty := toGenAIType(domainSchema.Type)
	properties := make(map[string]*genai.Schema)
	for name, prop := range domainSchema.Properties {
		properties[name] = structuredOutputSchemaToGenAI(&prop)
	}

	return &genai.Schema{
		Type:        ty,
		Description: domainSchema.Description,
		Properties:  properties,
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
