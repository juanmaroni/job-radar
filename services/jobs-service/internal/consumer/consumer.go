package consumer

import (
	"context"
	"encoding/json"
	"jobradar/events"
	"jobradar/jobs-service/internal/domain"
	"log/slog"
)

type Consumer struct {
    logger *slog.Logger
    reader kafkaReader
    store domain.JobStore
}

func NewConsumer(logger *slog.Logger, reader kafkaReader, store domain.JobStore) *Consumer {
	return &Consumer{
		logger: logger,
		reader: reader,
		store: store,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)

		if err != nil {
			return err
		}

		var jobMatched events.JobMatchedEvent

		if err = json.Unmarshal(msg.Value, &jobMatched); err != nil {
			c.logger.Error("consumer deserialization failed", slog.String("error", err.Error()))
			continue
		}

		job, err := domain.NewJob(jobMatched.Source, jobMatched.SourceId, jobMatched.Company, jobMatched.Role, jobMatched.Score)

		if err != nil {
			c.logger.Error("job construction failed", slog.String("error", err.Error()))
			continue
		}

		if err = c.store.Create(ctx, job); err != nil {
			c.logger.Error("job insert failed", slog.String("error", err.Error()))
			continue
		}
	}
}
