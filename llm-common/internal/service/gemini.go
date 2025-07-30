package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/compendium-tech/compendium/llm-common/internal/domain"
	"google.golang.org/genai"
	"strings"
	"time"
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

	return &geminiClient{
		client: client,
	}, nil
}

// AddTools registers tools with the Gemini client by converting domain ToolDefinitions.
func (c *geminiClient) AddTools(tools []domain.ToolDefinition) {
	var genAITools []*genai.Tool
	for _, tool := range tools {
		genAITool := &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  convertToolParameters(tool.Parameters),
				},
			},
		}

		genAITools = append(genAITools, genAITool)
	}

	c.tools = genAITools
}

// GenerateResponse generates a response from the Gemini model based on the provided messages.
func (c *GeminiClient) GenerateResponse(ctx context.Context, messages []domain.Message) (domain.Response, error) {
	chatSession := c.model.StartChat()
	var history []*genai.Content
	for _, msg := range messages {
		content := convertToGenaiContent(msg)
		history = append(history, content)
	}
	chatSession.History = history

	resp, err := chatSession.SendMessage(ctx, history[len(history)-1].Parts...)
	if err != nil {
		return domain.Response{}, fmt.Errorf("failed to send message: %w", err)
	}

	return convertToDomainResponse(resp)
}

// GenerateStructuredResponse generates a structured response using a JSON schema.
func (c *GeminiClient) GenerateStructuredResponse(ctx context.Context, messages []domain.Message, schema map[string]interface{}) (domain.Response, error) {
	chatSession := c.model.StartChat()
	var history []*genai.Content
	for _, msg := range messages {
		content := convertToGenaiContent(msg)
		history = append(history, content)
	}
	chatSession.History = history

	// Configure structured output with JSON schema
	c.model.GenerationConfig.ResponseMIMEType = "application/json"
	c.model.GenerationConfig.ResponseSchema = &genai.Schema{
		Type:       genai.TypeObject,
		Properties: convertSchemaToGenai(schema),
	}

	resp, err := chatSession.SendMessage(ctx, history[len(history)-1].Parts...)
	if err != nil {
		return domain.Response{}, fmt.Errorf("failed to send message: %w", err)
	}

	return convertToDomainResponse(resp)
}

// Helper functions for conversions between domain and Gemini types.

// convertToGenaiContent converts a domain Message to a Gemini Content.
func convertToGenaiContent(msg domain.Message) *genai.Content {
	content := &genai.Content{
		Role:  string(msg.Role),
		Parts: []*genai.Part{{Text: msg.TextContent}},
	}

	for _, toolCall := range msg.ToolCalls {
		params, _ := json.Marshal(toolCall.Parameters)
		content.Parts = append(content.Parts, genai.Part{
			FunctionCall: &genai.FunctionCall{
				Name: toolCall.Name,
				Args: json.RawMessage(params),
			},
		})
	}

	return content
}

// convertToDomainResponse converts a Gemini response to a domain Response.
func convertToDomainResponse(resp *genai.GenerateContentResponse) (domain.Response, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return domain.Response{}, fmt.Errorf("no content in response")
	}

	message := domain.Message{
		Role: domain.ROLE_ASSISTANT,
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		switch {
		case part.Text != "":
			message.TextContent = part.Text
		case part.FunctionCall != nil:
			var params map[string]interface{}
			if err := json.Unmarshal(part.FunctionCall.Args, &params); err != nil {
				return domain.Response{}, fmt.Errorf("failed to unmarshal function call args: %w", err)
			}
			message.ToolCalls = append(message.ToolCalls, domain.ToolCall{
				ID:         generateToolCallID(part.FunctionCall.Name),
				Name:       part.FunctionCall.Name,
				Parameters: params,
			})
		}
	}

	return domain.Response{Message: message}, nil
}

// convertToolParameters converts domain ToolParameters to Gemini Schema.
func convertToolParameters(params []domain.ToolParameter) *genai.Schema {
	properties := make(map[string]*genai.Schema)
	var required []string

	for _, param := range params {
		schema := &genai.Schema{
			Type:        genai.Type(strings.ToLower(param.Type)),
			Description: param.Description,
		}
		if len(param.Enum) > 0 {
			schema.Enum = param.Enum
		}
		properties[param.Name] = schema
		if param.IsRequired {
			required = append(required, param.Name)
		}
	}

	return &genai.Schema{
		Type:       genai.TypeObject,
		Properties: properties,
		Required:   required,
	}
}

func convertSchemaToGenAI(schema map[string]interface{}) map[string]*genai.Schema {
	properties := make(map[string]*genai.Schema)
	for key, val := range schema {
		schemaType := inferSchemaType(val)
		propSchema := &genai.Schema{Type: schemaType}
		if schemaType == genai.TypeObject {
			if nested, ok := val.(map[string]interface{}); ok {
				propSchema.Properties = convertSchemaToGenAI(nested)
			}
		} else if schemaType == genai.TypeArray {
			if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
				propSchema.Items = &genai.Schema{Type: inferSchemaType(arr[0])}
			}
		}

		properties[key] = propSchema
	}
	return properties
}

func inferSchemaType(val interface{}) genai.Type {
	switch val.(type) {
	case map[string]interface{}:
		return genai.TypeObject
	case []interface{}:
		return genai.TypeArray
	case string:
		return genai.TypeString
	case float64, int, int64:
		return genai.TypeNumber
	case bool:
		return genai.TypeBoolean
	default:
		return genai.TypeString
	}
}

func generateToolCallID(name string) string {
	return fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
}
