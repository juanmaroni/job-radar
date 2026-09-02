package main

import (
	"context"
	"errors"
	"fmt"
	"jobradar/tracker-service/internal/config/environment"
	"jobradar/tracker-service/internal/config/logging"
	"jobradar/tracker-service/internal/config/rdb"
	"jobradar/tracker-service/internal/config/server"
	"jobradar/tracker-service/internal/domain"
	"jobradar/tracker-service/internal/handler/application"
	"jobradar/tracker-service/internal/handler/health"
	"jobradar/tracker-service/internal/repository"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Tracker service started.")

	// PostgreSQL setup
	dbEnv, err := environment.GetEnvVar("TRACKER_SERVICE_DATABASE_URI")

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

	// Redis setup
	redisClient, err := rdb.CreateClient()

	if err != nil {
		logger.Error("redis client creation failed, continuing without cache", slog.String("error", err.Error()))
	}

	repo := repository.NewApplicationRepository(logger, pool)

	var appStore domain.ApplicationStore = repo
		
	if redisClient != nil {
		defer redisClient.Close()
		ttl, err := rdb.LoadTtl()

		if err != nil {
			logger.Warn("redis ttl load, using default ttl", slog.String("error", err.Error()))
		}

		appStore = repository.NewCachedApplicationRepository(logger, repo, redisClient, ttl)
	}

	appHandler := application.NewApplicationHandler(logger, appStore)

	// Server setup
	srv, err := server.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	serverUri := srv.Uri
	mux := http.NewServeMux()
	mux.Handle("GET /health", health.NewHealthHandler(logger))
	mux.HandleFunc("POST /applications", appHandler.CreateApplication)
	mux.HandleFunc("GET /applications", appHandler.ListApplications)

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
