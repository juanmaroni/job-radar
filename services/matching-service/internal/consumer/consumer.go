package consumer

import (
	"context"
	"encoding/json"
	"jobradar/events"
	"jobradar/matching-service/internal/scoring"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	logger *slog.Logger
	reader kafkaReader
	writer kafkaWriter
	keywords []string
}

func NewConsumer(logger *slog.Logger, reader kafkaReader, writer kafkaWriter, keywords []string) *Consumer {
	return &Consumer{
		logger: logger,
		reader: reader,
		writer: writer,
		keywords: keywords,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)

		if err != nil {
			return err
		}

		var jobPublished events.JobPublishedEvent
		
		if err = json.Unmarshal(msg.Value, &jobPublished); err != nil {
			c.logger.Error("consumer deserialization failed", slog.String("error", err.Error()))
			continue
		}

		score := scoring.CalculateScore(jobPublished.Description, c.keywords)

		jobMatched := events.JobMatchedEvent{
			SourceId: jobPublished.SourceId,
			Source: jobPublished.Source,
			Company: jobPublished.Company,
			Role: jobPublished.Role,
			OfferUrl: jobPublished.OfferUrl,
			Location: jobPublished.Location,
			SalaryMin: jobPublished.SalaryMin,
			SalaryMax: jobPublished.SalaryMax,
			Description: jobPublished.Description,
			Score: score,
			PublishedAt: jobPublished.PublishedAt,
		}

		jsonJobMatched, err := json.Marshal(&jobMatched)

		if err != nil {
			c.logger.Error("consumer serialization failed", slog.String("error", err.Error()))
			continue
		}

		if err = c.writer.WriteMessages(ctx, kafka.Message{
			Key: []byte(jobPublished.SourceId),
			Value: jsonJobMatched,
		}); err != nil {
			c.logger.Error("consumer writing failed", slog.String("error", err.Error()))
			continue
		}
	}
}
