package server

import "jobradar/matching-service/internal/config/environment"

type Server struct {
	Port string
}

func Load() (*Server, error) {
	serverPort, err := environment.GetEnvVar("MATCHING_SERVICE_GRPC_PORT")

	if err != nil {
		return nil, err
	}

	return &Server {
		Port: serverPort,
	}, nil
}
