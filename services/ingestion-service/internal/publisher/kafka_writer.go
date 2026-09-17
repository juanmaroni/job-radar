package publisher

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaWriter interface {
    WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}
