package main

import (
	"context"
	"errors"
	"fmt"
	"jobradar/jobs-service/internal/config/broker"
	"jobradar/jobs-service/internal/config/environment"
	"jobradar/jobs-service/internal/config/logging"
	"jobradar/jobs-service/internal/config/server"
	"jobradar/jobs-service/internal/consumer"
	"jobradar/jobs-service/internal/handler/health"
	"jobradar/jobs-service/internal/handler/job"
	"jobradar/jobs-service/internal/repository"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Jobs service started.")

	// PostgreSQL setup
	dbEnv, err := environment.GetEnvVar("JOBS_SERVICE_DATABASE_URI")

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbEnv)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// Kafka setup
	brokerCfg, err := broker.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerCfg.Uri},
		Topic: "jobs.matched",
		GroupID: "jobs-service",
	})

	repo := repository.NewJobRepository(logger, pool)
	consumer := consumer.NewConsumer(logger, kafkaReader, repo)

	go func() {
		logger.Info(fmt.Sprintf("Starting consumer at '%s'.", brokerCfg.Uri))

		if err := consumer.Run(ctx); err != nil {
			logger.Error("consumer stopped", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	jobHandler := job.NewJobHandler(logger, repo)

	// Server setup
	srv, err := server.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	serverUri := srv.Uri
	mux := http.NewServeMux()
	mux.Handle("GET /health", health.NewHealthHandler(logger))
	mux.HandleFunc("GET /jobs", jobHandler.ListJobs)

	logger.Info(fmt.Sprintf("Starting server at '%s'.", serverUri))

	logMux := logging.HttpRequestLogging(logger)(mux)
	err = http.ListenAndServe(serverUri, logMux)

	if errors.Is(err, http.ErrServerClosed) {
		logger.Error("Error: server closed.")
		os.Exit(1)
	} else if err != nil {
		logger.Error(fmt.Sprintf("Error: couldn't start server, %s\n", err))
		os.Exit(1)
	}
}
