package grpcv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/anypb"

	pbhelp "github.com/compendium-tech/compendium/common/pkg/pb"

	"github.com/compendium-tech/compendium/llm-service/internal/domain"
	pb "github.com/compendium-tech/compendium/llm-service/internal/proto/v1"
	"github.com/compendium-tech/compendium/llm-service/internal/service"
)

type LLMServiceServer struct {
	pb.UnimplementedLLMServiceServer
	llmService service.LLMService
}

func NewLLMServiceServer(llmService service.LLMService) LLMServiceServer {
	return LLMServiceServer{llmService: llmService}
}

func (s LLMServiceServer) Register(server *grpc.Server) {
	pb.RegisterLLMServiceServer(server, s)
	reflection.Register(server)
}

func (s LLMServiceServer) GenerateResponse(ctx context.Context, req *pb.GenerateResponseRequest) (*pb.GenerateResponseResponse, error) {
	chatHistory := make([]domain.Message, len(req.ChatHistory))
	for i, msg := range req.ChatHistory {
		toolCalls := make([]domain.ToolCall, len(msg.ToolCalls))
		for j, tc := range msg.ToolCalls {
			var err error

			params := make(map[string]any)
			for k, v := range tc.Parameters {
				params[k], err = pbhelp.AnyPBToAny(v)
				if err != nil {
					return nil, err
				}
			}

			toolCalls[j] = domain.ToolCall{
				ID:         tc.Id,
				Name:       tc.Name,
				Parameters: params,
			}
		}

		chatHistory[i] = domain.Message{
			Role:      rolePBToRole(msg.Role),
			Text:      msg.Text,
			ToolCalls: toolCalls,
		}
	}

	tools := make([]domain.ToolDefinition, len(req.Tools))
	for i, tool := range req.Tools {
		tools[i] = domain.ToolDefinition{
			Name:             tool.Name,
			Description:      tool.Description,
			ParametersSchema: schemaPBToSchema(tool.ParametersSchema),
		}
	}

	var schema *domain.Schema
	if req.StructuredOutputSchema != nil {
		schema = schemaPBToSchema(req.StructuredOutputSchema)
	}

	resp, err := s.llmService.GenerateResponse(ctx, chatHistory, tools, schema)
	if err != nil {
		return nil, err
	}

	toolCalls := make([]*pb.ToolCall, len(resp.ToolCalls))
	for i, tc := range resp.ToolCalls {
		params := make(map[string]*anypb.Any)
		for k, v := range tc.Parameters {
			params[k], err = pbhelp.AnyToAnyPB(v)
			if err != nil {
				return nil, err
			}
		}

		toolCalls[i] = &pb.ToolCall{
			Id:         tc.ID,
			Name:       tc.Name,
			Parameters: params,
		}
	}

	return &pb.GenerateResponseResponse{
		Message: &pb.Message{
			Role:      roleToRolePB(resp.Role),
			Text:      resp.Text,
			ToolCalls: toolCalls,
		},
	}, nil
}

func roleToRolePB(role domain.Role) pb.Role {
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

func rolePBToRole(role pb.Role) domain.Role {
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

func schemaPBToSchema(protoSchema *pb.Schema) *domain.Schema {
	if protoSchema == nil {
		return nil
	}

	properties := make(map[string]domain.Schema)
	for k, v := range protoSchema.Properties {
		properties[k] = *schemaPBToSchema(v)
	}

	var items *domain.Schema
	if protoSchema.Items != nil {
		items = schemaPBToSchema(protoSchema.Items)
	}

	return &domain.Schema{
		Type:        domain.Type(protoSchema.Type.Value),
		Description: protoSchema.Description,
		Properties:  properties,
		Items:       items,
		MaxItems:    &protoSchema.MaxItems,
		MinItems:    &protoSchema.MinItems,
		Required:    protoSchema.Required,
	}
}
