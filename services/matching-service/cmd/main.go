package main

import (
	"fmt"
	"jobradar/matching-service/internal/config/logging"
	configserver "jobradar/matching-service/internal/config/server"
	matchingserver "jobradar/matching-service/internal/server"
	matching "jobradar/proto/gen/go/matching/v1"
	"net"

	"google.golang.org/grpc"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Matching service started.")

	// gRPC server setup
	srv, err := configserver.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	serverPort := srv.Port

	lis, err := net.Listen("tcp", serverPort)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	logger.Info(fmt.Sprintf("Starting server at '%s'.", serverPort))

	grpcServer := grpc.NewServer()

	matching.RegisterMatchingServiceServer(grpcServer, matchingserver.NewMatchingServer())

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
