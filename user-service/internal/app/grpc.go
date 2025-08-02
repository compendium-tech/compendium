package app

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/log"
	netapp "github.com/compendium-tech/compendium/common/pkg/net"

	"github.com/compendium-tech/compendium/user-service/internal/config"
	grpcv1 "github.com/compendium-tech/compendium/user-service/internal/delivery/grpc/v1"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/service"
)

type GrpcAppDependencies struct {
	Config config.GrpcAppConfig
	PgDB   *sql.DB
}

func NewGrpcApp(deps GrpcAppDependencies) netapp.GrpcApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "user-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	userRepository := repository.NewPgUserRepository(deps.PgDB)
	userService := service.NewUserService(userRepository)

	grpcServer := grpc.NewServer()
	grpcv1.NewUserServiceServer(userService).Register(grpcServer)

	return netapp.NewGrpcApp(grpcServer)
}
