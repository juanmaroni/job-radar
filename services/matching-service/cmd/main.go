package main

import (
	"fmt"
	"jobradar/matching-service/internal/config/logging"
	configserver "jobradar/matching-service/internal/config/server"
	matchingserver "jobradar/matching-service/internal/server"
	matchingproto "jobradar/proto/gen/go/matching/v1"
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

	lis, err := net.Listen("tcp", srv.Port)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	logger.Info(fmt.Sprintf("Starting server at '%s'.", srv.Port))

	grpcServer := grpc.NewServer()

	matchingproto.RegisterMatchingServiceServer(grpcServer, matchingserver.NewMatchingServer())

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
