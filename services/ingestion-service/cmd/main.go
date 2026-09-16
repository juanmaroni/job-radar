package main

import (
	"context"
	"fmt"
	"jobradar/ingestion-service/internal/config/broker"
	"jobradar/ingestion-service/internal/config/logging"
	"jobradar/ingestion-service/internal/publisher"
	"jobradar/ingestion-service/internal/source"
	"log/slog"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	logger := logging.LoadJSONLogger()
	logger.Info("Ingestion service started.")

	// Kafka setup
	brokerCfg, err := broker.Load()

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	kafkaWriter := &kafka.Writer{
		Addr: kafka.TCP(brokerCfg.Uri),
		Topic: "jobs.new",
		AllowAutoTopicCreation: true,
	}
	httpCli := &http.Client{
		Timeout: 5 * time.Second,
	}
	remoteokSrc := source.NewRemoteOKSource(logger, httpCli, source.RemoteOkApiUrl)
	kafkaJobPublisher := publisher.NewKafkaJobPublisher(logger, kafkaWriter)
	ctx := context.Background()
	ticker, err := broker.LoadTicker()

	if err != nil {
		logger.Warn("ticker load, using default interval", slog.String("error", err.Error()))
	}
	
	defer ticker.Stop()

	logger.Info(fmt.Sprintf("Publishing jobs to broker at '%s'.", brokerCfg.Uri))

	for {
		select {
		case <-ticker.C:
			jobs, err := remoteokSrc.FetchJobs(ctx)

			if err != nil {
				logger.Error("failed to fetch jobs", slog.String("error", err.Error()))
				continue
			}

			var numPublished int = 0
			var numFailed int = 0

			for _, job := range jobs {
				if err = kafkaJobPublisher.Publish(ctx, job); err != nil {
					numFailed += 1
				} else {
					numPublished += 1
				}
			}

			logger.Info("ingestion cycle completed", slog.Int("published", numPublished), slog.Int("failed", numFailed))
		case <-ctx.Done():
			return
		}
	}
}
