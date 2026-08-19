package main

import (
	"errors"
	"fmt"
	"jobradar/tracker-service/internal/config/logging"
	"jobradar/tracker-service/internal/config/server"
	"jobradar/tracker-service/internal/handler/health"
	"net/http"
	"os"
)

func main() {
	// Logger
	logger := logging.LoadJSONLogger()
	logger.Info("Application started.")
	
	// Load environment vars
	srv, err := server.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	logger.Info("Environment variables loaded.")

	// Server set up
	serverUri := srv.Uri
	mux := http.NewServeMux()
	mux.Handle("/health", health.NewHealthHandler(logger))

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
