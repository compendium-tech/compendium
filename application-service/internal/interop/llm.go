package interop

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/compendium-tech/compendium/application-service/internal/domain"
	pb "github.com/compendium-tech/compendium/application-service/internal/proto/v1"
	pbhelp "github.com/compendium-tech/compendium/common/pkg/pb"
)

type llmServiceGrpcClient struct {
	client pb.LLMServiceClient
}

type LLMService interface {
	GenerateResponse(
		ctx context.Context, chatHistory []domain.LLMMessage,
		tools []domain.LLMToolDefinition, structuredOutputSchema *domain.LLMSchema) domain.LLMMessage
}

func NewGrpcLLMServiceClient(target string) (LLMService, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	c := pb.NewLLMServiceClient(conn)

	return &llmServiceGrpcClient{
		client: c,
	}, nil
}

func (c *llmServiceGrpcClient) GenerateResponse(
	ctx context.Context,
	chatHistory []domain.LLMMessage,
	tools []domain.LLMToolDefinition,
	structuredOutputSchema *domain.LLMSchema,
) domain.LLMMessage {
	protoChatHistory := make([]*pb.Message, len(chatHistory))

	for i, msg := range chatHistory {
		toolCalls := make([]*pb.ToolCall, len(msg.ToolCalls))
		for j, tc := range msg.ToolCalls {
			var err error

			params := make(map[string]*anypb.Any)
			for k, v := range tc.Parameters {
				params[k], err = pbhelp.AnyToAnyPB(v)
				if err != nil {
					panic(err)
				}
			}

			toolCalls[j] = &pb.ToolCall{
				Id:         tc.ID,
				Name:       tc.Name,
				Parameters: params,
			}
		}

		protoChatHistory[i] = &pb.Message{
			Role:      roleToRolePB(msg.Role),
			Text:      msg.Text,
			ToolCalls: toolCalls,
		}
	}

	protoTools := make([]*pb.ToolDefinition, len(tools))
	for i, tool := range tools {
		protoTools[i] = &pb.ToolDefinition{
			Name:             tool.Name,
			Description:      tool.Description,
			ParametersSchema: schemaToSchemaPB(tool.ParametersSchema),
		}
	}

	var protoSchema *pb.Schema
	if structuredOutputSchema != nil {
		protoSchema = schemaToSchemaPB(structuredOutputSchema)
	}

	resp, err := c.client.GenerateResponse(ctx, &pb.GenerateResponseRequest{
		ChatHistory:            protoChatHistory,
		Tools:                  protoTools,
		StructuredOutputSchema: protoSchema,
	})
	if err != nil {
		panic(err)
	}

	toolCalls := make([]domain.LLMToolCall, len(resp.Message.ToolCalls))
	for i, tc := range resp.Message.ToolCalls {
		params := make(map[string]any)
		for k, v := range tc.Parameters {
			params[k], err = pbhelp.AnyPBToAny(v)
			if err != nil {
				panic(err)
			}
		}

		toolCalls[i] = domain.LLMToolCall{
			ID:         tc.Id,
			Name:       tc.Name,
			Parameters: params,
		}
	}

	return domain.LLMMessage{
		Role:      rolePBToRole(resp.Message.Role),
		Text:      resp.Message.Text,
		ToolCalls: toolCalls,
	}
}

func roleToRolePB(role domain.LLMRole) pb.Role {
	switch role {
	case domain.RoleSystem:
		return pb.Role_SYSTEM
	case domain.RoleUser:
		return pb.Role_USER
	case domain.RoleAssistant:
		return pb.Role_ASSISTANT
	default:
		return pb.Role_SYSTEM
	}
}

func rolePBToRole(role pb.Role) domain.LLMRole {
	switch role {
	case pb.Role_SYSTEM:
		return domain.RoleSystem
	case pb.Role_USER:
		return domain.RoleUser
	case pb.Role_ASSISTANT:
		return domain.RoleAssistant
	default:
		return domain.RoleSystem
	}
}

func schemaToSchemaPB(domainSchema *domain.LLMSchema) *pb.Schema {
	if domainSchema == nil {
		return nil
	}

	properties := make(map[string]*pb.Schema)
	for k, v := range domainSchema.Properties {
		properties[k] = schemaToSchemaPB(&v)
	}

	var items *pb.Schema
	if domainSchema.Items != nil {
		items = schemaToSchemaPB(domainSchema.Items)
	}

	return &pb.Schema{
		Type:        &pb.Type{Value: string(domainSchema.Type)},
		Description: domainSchema.Description,
		Properties:  properties,
		Items:       items,
		MaxItems:    *domainSchema.MaxItems,
		MinItems:    *domainSchema.MinItems,
		Required:    domainSchema.Required,
	}
}
