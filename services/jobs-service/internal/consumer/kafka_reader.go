package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}
