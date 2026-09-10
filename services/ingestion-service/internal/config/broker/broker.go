package broker

import (
	"jobradar/ingestion-service/internal/config/environment"
	"strconv"
	"time"
)

const defaultIntervalSeconds = 300

type Broker struct {
	Uri string
}

func Load() (*Broker, error) {
	kafkaUri, err := environment.GetEnvVar("INGESTION_SERVICE_KAFKA_BROKER")

	if err != nil {
		return nil, err
	}

	return &Broker {
		Uri: kafkaUri,
	}, nil
}

func LoadTicker() (*time.Ticker, error) {
	kafkaTickerStr, err := environment.GetEnvVar("INGESTION_SERVICE_INTERVAL_SECONDS")

	if err != nil {
		return time.NewTicker(defaultIntervalSeconds * time.Second), err
	}

	seconds, err := strconv.Atoi(kafkaTickerStr)

	if err != nil {
		return time.NewTicker(defaultIntervalSeconds * time.Second), err
	} else {
		return time.NewTicker(time.Duration(seconds) * time.Second), nil
	}
}
