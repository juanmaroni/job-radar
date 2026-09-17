package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"jobradar/events"
	"jobradar/ingestion-service/internal/domain"
	"log/slog"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
    lastMessages []kafka.Message
    err          error
}

func (f *fakeWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
    if f.err != nil {
        return f.err
    }

    f.lastMessages = msgs

    return nil
}

func TestPublish_EventCreated(t *testing.T) {
	ctx := context.Background()
	writer := &fakeWriter{
		lastMessages: []kafka.Message{},
		err: nil,
	}

	jobPublisher := NewKafkaJobPublisher(
		slog.New(slog.DiscardHandler),
		writer,
	)

	job, err := domain.NewJob("RemoteOK", "1", "Acme", "Backend Engineer", "https://www.acme.com/jobs/1", "Remote", "Go and PostgreSQL", nil, nil, nil)

	require.NoError(t, err)

	err = jobPublisher.Publish(ctx, *job)

	require.NoError(t, err)
	require.Len(t, writer.lastMessages, 1)

	var event events.JobPublishedEvent
    err = json.Unmarshal(writer.lastMessages[0].Value, &event)
    require.NoError(t, err)

    assert.Equal(t, "Acme", event.Company)
    assert.Equal(t, "Backend Engineer", event.Role)
	assert.Equal(t, "https://www.acme.com/jobs/1", event.OfferUrl)
	assert.Equal(t, "Go and PostgreSQL", event.Description)
}

func TestPublish_WriteMessagesError(t *testing.T) {
	ctx := context.Background()
	writer := &fakeWriter{
		lastMessages: []kafka.Message{},
		err: errors.New("write messages failure"),
	}

	jobPublisher := NewKafkaJobPublisher(
		slog.New(slog.DiscardHandler),
		writer,
	)

	job, err := domain.NewJob("RemoteOK", "1", "Acme", "Backend Engineer", "https://www.acme.com/jobs/1", "Remote", "Go and PostgreSQL", nil, nil, nil)

	require.NoError(t, err)

	err = jobPublisher.Publish(ctx, *job)

	require.Error(t, err)
}
