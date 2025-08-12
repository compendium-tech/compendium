package netapp

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GrpcApp struct {
	server *grpc.Server
}

func NewGrpcApp(server *grpc.Server) GrpcApp {
	return GrpcApp{
		server: server,
	}
}

func (a GrpcApp) Run() error {
	port, err := getPortEnv()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logrus.Infof("Starting gRPC server on :%d", port)
	return a.server.Serve(lis)
}
