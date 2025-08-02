package app

import (
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/log"
	netapp "github.com/compendium-tech/compendium/common/pkg/net"

	"github.com/compendium-tech/compendium/llm-service/internal/config"
	grpcv1 "github.com/compendium-tech/compendium/llm-service/internal/delivery/grpc/v1"
	"github.com/compendium-tech/compendium/llm-service/internal/service"
)

type Dependencies struct {
	Config     *config.AppConfig
	LLMService service.LLMService
}

func NewApp(deps Dependencies) netapp.GrpcApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "llm-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	grpcServer := grpc.NewServer()
	grpcv1.NewLLMServiceServer(deps.LLMService).Register(grpcServer)

	return netapp.NewGrpcApp(grpcServer)
}
