package main

import (
	"context"
	"errors"
	"fmt"
	"jobradar/tracker-service/internal/config/environment"
	"jobradar/tracker-service/internal/config/logging"
	"jobradar/tracker-service/internal/config/server"
	"jobradar/tracker-service/internal/handler/application"
	"jobradar/tracker-service/internal/handler/health"
	"jobradar/tracker-service/internal/repository"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Application started.")
	
	srv, err := server.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

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

	if err = pool.Ping(ctx); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	repo := repository.NewApplicationRepository(pool)
	appHandler := application.NewApplicationHandler(logger, repo)

	// Server set up
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
