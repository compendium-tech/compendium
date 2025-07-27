package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/compendium-tech/compendium/assistant-service/internal/proto/v1"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genai"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type AssistantServiceServer struct {
	pb.UnimplementedAssistantServiceServer
	geminiClient *genai.Client
}

func NewAssistantServiceServer(ctx context.Context) (*AssistantServiceServer, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	return &AssistantServiceServer{
		geminiClient: client,
	}, nil
}

func (s *AssistantServiceServer) CreateResponse(ctx context.Context, req *pb.CreateResponseRequest) (*pb.CreateResponseResponse, error) {
	var chatMessages []*genai.Content

	for _, msg := range req.Messages {
		role := convertRoleToGemini(msg.Role)

		chatMessages = append(chatMessages, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: msg.TextContent}},
		})
	}

	var tools []*genai.Tool

	for _, toolDef := range req.Tools {
		var requiredParameters []string
		var parameters = make(map[string]*genai.Schema)

		for _, param := range toolDef.Parameters {
			parameters[param.Name] = &genai.Schema{
				Type:        convertTypeToGemini(param.Type),
				Description: param.Description,
				Enum:        param.Enum,
			}

			if param.IsRequired {
				requiredParameters = append(requiredParameters, param.Name)
			}
		}

		tools = append(tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{{
				Name:        toolDef.Name,
				Description: toolDef.Description,
				Parameters: &genai.Schema{
					Type:       genai.TypeObject,
					Properties: parameters,
				},
			}},
		})
	}

	resp, err := s.geminiClient.Models.GenerateContent(ctx, "gemini-2.0-flash", chatMessages, &genai.GenerateContentConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return &pb.CreateResponseResponse{
			Message: nil,
		}, nil
	}

	candidate := resp.Candidates[0]

	var texts []string
	var functionCalls []*genai.FunctionCall
	for _, part := range candidate.Content.Parts {
		texts = append(texts, part.Text)

		if part.FunctionCall != nil {
			functionCalls = append(functionCalls, part.FunctionCall)
		}
	}

	responseMessage := &pb.Message{
		Role:        pb.Message_ASSISTANT,
		TextContent: strings.Join(texts, ""),
	}

	if len(functionCalls) > 0 {
		var toolCalls []*pb.ToolCall
		for _, fc := range functionCalls {
			toolCall := &pb.ToolCall{
				Id:   fc.ID,
				Name: fc.Name,
			}
			for paramName, paramValue := range fc.Args {
				anyVal, err := convertInterfaceToAny(paramValue)
				if err != nil {
					return nil, fmt.Errorf("failed to convert parameter to Any: %v", err)
				}

				toolCall.Parameters = append(toolCall.Parameters, &pb.ToolCallParameter{
					Name:  paramName,
					Value: anyVal,
				})
			}
			toolCalls = append(toolCalls, toolCall)
		}

		responseMessage.ToolCalls = toolCalls
	}

	return &pb.CreateResponseResponse{
		Message: responseMessage,
	}, nil
}

func convertRoleToGemini(role pb.Message_Role) string {
	switch role {
	case pb.Message_SYSTEM:
		return "system"
	case pb.Message_USER:
		return "user"
	case pb.Message_ASSISTANT:
		return "assistant"
	case pb.Message_TOOL:
		return "tool"
	default:
		return "user"
	}
}

func convertTypeToGemini(ty string) genai.Type {
	switch ty {
	case "string":
		return genai.TypeString
	case "number":
		return genai.TypeNumber
	case "array":
		return genai.TypeArray
	case "object":
		return genai.TypeObject
	default:
		return genai.TypeString
	}
}

func convertInterfaceToAny(v interface{}) (*anypb.Any, error) {
	anyValue := &anypb.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}
	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	return anyValue, err
}

func main() {
	ctx := context.Background()

	server, err := NewAssistantServiceServer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAssistantServiceServer(grpcServer, server)
	log.Printf("Server listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
