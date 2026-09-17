package publisher

import (
	"context"
	"encoding/json"
	"jobradar/events"
	"jobradar/ingestion-service/internal/domain"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type KafkaJobPublisher struct {
	logger *slog.Logger
	writer kafkaWriter
}

func NewKafkaJobPublisher(logger *slog.Logger, writer kafkaWriter) *KafkaJobPublisher {
    return &KafkaJobPublisher{
		logger: logger,
		writer: writer,
	}
}

func (kjp *KafkaJobPublisher) Publish(ctx context.Context, job domain.Job) error {
	jobPublished := events.JobPublishedEvent {
		SourceId: job.SourceId,
		Source: job.Source,
		Company: job.Company,
		Role: job.Role,
		OfferUrl: job.OfferUrl,
		Description: job.Description,
	}

	jsonJobPublished, err := json.Marshal(jobPublished)

	if err != nil {
		kjp.logger.Error("json marshal fail", slog.String("error", err.Error()))

		return err
	}

	if err = kjp.writer.WriteMessages(ctx, kafka.Message{
		Key: []byte(job.SourceId),
		Value: jsonJobPublished},
	); err != nil {
		kjp.logger.Error("kafka writer fail", slog.String("error", err.Error()))

		return err
	}

	return nil
}
