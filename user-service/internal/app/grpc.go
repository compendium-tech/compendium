package app

import (
	"database/sql"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/log"

	"github.com/compendium-tech/compendium/user-service/internal/config"
	grpcv1 "github.com/compendium-tech/compendium/user-service/internal/delivery/grpc/v1"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/service"
)

type GrpcApp struct {
	server *grpc.Server
	port   uint16
}

func (g *GrpcApp) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	logrus.Infof("Starting gRPC server on :%d", g.port)
	return g.server.Serve(lis)
}

type GrpcAppDependencies struct {
	Config config.GrpcAppConfig
	PgDB   *sql.DB
}

func NewGrpcApp(deps GrpcAppDependencies) *GrpcApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	userRepository := repository.NewPgUserRepository(deps.PgDB)
	userService := service.NewUserService(userRepository)

	grpcServer := grpc.NewServer()
	grpcv1.NewUserServiceServer(userService).Register(grpcServer)

	return &GrpcApp{server: grpcServer, port: deps.Config.GrpcPort}
}
