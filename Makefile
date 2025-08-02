protoc:
	protoc --go_out=user-service --go-grpc_out=user-service user-service/proto/user_service.proto
	protoc --go_out=subscription-service --go-grpc_out=subscription-service subscription-service/proto/user_service.proto
	protoc --go_out=llm-service --go-grpc_out=llm-service llm-service/proto/llm_service.proto
	protoc --go_out=application-service --go-grpc_out=application-service application-service/proto/llm_service.proto
