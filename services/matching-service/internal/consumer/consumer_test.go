package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"jobradar/events"
	"log/slog"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
    lastMessages []kafka.Message
    err error
}

func (fw *fakeWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
    if fw.err != nil {
        return fw.err
    }

    fw.lastMessages = msgs

    return nil
}

type fakeReader struct {
	messages []kafka.Message
    err error
	calls int
}

func (fr *fakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if fr.calls >= len(fr.messages) {
        if fr.err != nil {
            return kafka.Message{}, fr.err
        }
		
        return kafka.Message{}, errors.New("no more messages")
    }

    msg := fr.messages[fr.calls]
    fr.calls += 1

    return msg, nil
}

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	reader := &fakeReader{
		messages: []kafka.Message{
			{Value: []byte(`{"source": "RemoteOK", "source_id": "1", "company": "Acme", "role": "Backend Go Developer", "description": "Go and PostgreSQL"}`)},
		},
		err: errors.New("no more messages"),
	}
	writer := &fakeWriter{
		lastMessages: []kafka.Message{},
		err: nil,
	}

	consumer := NewConsumer(slog.New(slog.DiscardHandler), reader, writer, []string{"Go"})
	err := consumer.Run(ctx)

	require.Error(t, err)
	assert.ErrorContains(t, err, "no more messages")
	assert.Equal(t, reader.calls, 1)

	require.Len(t, writer.lastMessages, 1)

	var matched events.JobMatchedEvent
	err = json.Unmarshal(writer.lastMessages[0].Value, &matched)
	require.NoError(t, err)

	assert.Equal(t, "Acme", matched.Company)
	assert.Equal(t, int32(1), matched.Score)
}
