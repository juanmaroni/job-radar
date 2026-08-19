package server

import "jobradar/tracker-service/internal/config/environment"

type Server struct {
	Uri string
}

func Load() (*Server, error) {
	serverUri, err := environment.GetEnvVar("TRACKER_SERVICE_URI")

	if err != nil {
		return nil, err
	}

	return &Server {
		Uri: serverUri,
	}, nil
}
