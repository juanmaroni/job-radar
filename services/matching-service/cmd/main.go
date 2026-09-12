package main

import (
	"context"
	"fmt"
	"jobradar/matching-service/internal/config/broker"
	"jobradar/matching-service/internal/config/logging"
	configserver "jobradar/matching-service/internal/config/server"
	"jobradar/matching-service/internal/consumer"
	matchingserver "jobradar/matching-service/internal/server"
	matchingproto "jobradar/proto/gen/go/matching/v1"
	"log/slog"
	"net"

	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Matching service started.")

	// Kafka setup
	brokerCfg, err := broker.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerCfg.Uri},
		Topic: "jobs.new",
		GroupID: "matching-service",
	})

	kafkaWriter := &kafka.Writer{
		Addr: kafka.TCP(brokerCfg.Uri),
		Topic: "jobs.matched",
		AllowAutoTopicCreation: true,
	}

	consumer := consumer.NewConsumer(logger, kafkaReader, kafkaWriter)
	ctx := context.Background()

	go func() {
		logger.Info(fmt.Sprintf("Starting consumer at '%s'.", brokerCfg.Uri))

		if err := consumer.Run(ctx); err != nil {
			logger.Error("consumer stopped", slog.String("error", err.Error()))
			panic(err)
		}
	}()

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
