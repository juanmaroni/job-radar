package broker

import (
	"jobradar/jobs-service/internal/config/environment"
)

const defaultIntervalSeconds = 300

type Broker struct {
	Uri string
}

func Load() (*Broker, error) {
	kafkaUri, err := environment.GetEnvVar("JOBS_SERVICE_KAFKA_BROKER")

	if err != nil {
		return nil, err
	}

	return &Broker {
		Uri: kafkaUri,
	}, nil
}
