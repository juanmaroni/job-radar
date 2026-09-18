package consumer

import (
	"context"
	"errors"
	"jobradar/jobs-service/internal/domain"
	"log/slog"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

type fakeJobStore struct {
	createErr error
    listErr error
    jobs []domain.Job
}

func (f *fakeJobStore) Create(ctx context.Context, job *domain.Job) error {
	if f.createErr != nil {
		return f.createErr
	}

	job.Id = int64(len(f.jobs) + 1)
	f.jobs = append(f.jobs, *job)

	return nil
}

func (f *fakeJobStore) List(ctx context.Context) ([]domain.Job, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}

	return f.jobs, nil
}

func TestRun_Success(t *testing.T) {
	ctx := context.Background()
	reader := &fakeReader{
		messages: []kafka.Message{
			{Value: []byte(`{"source": "RemoteOK", "source_id": "1", "company": "Acme", "role": "Backend Go Developer", "score": 1}`)},
		},
		err: errors.New("no more messages"),
	}
	jobStore := &fakeJobStore{
		createErr: nil,
		listErr: nil,
		jobs: []domain.Job{},
	}
	consumer := NewConsumer(
		slog.New(slog.DiscardHandler),
		reader,
		jobStore,
	)

	err := consumer.Run(ctx)

	require.Error(t, err)
	assert.ErrorContains(t, err, "no more messages")
	assert.Equal(t, reader.calls, 1)
	require.Len(t, jobStore.jobs, 1)
	assert.Equal(t, "RemoteOK", jobStore.jobs[0].Source)
	assert.Equal(t, "1", jobStore.jobs[0].SourceId)
	assert.Equal(t, "Acme", jobStore.jobs[0].Company)
	assert.Equal(t, "Backend Go Developer", jobStore.jobs[0].Role)
	assert.Equal(t, int32(1), jobStore.jobs[0].Score)
}

func TestRun_FailUnmarshal(t *testing.T) {
	ctx := context.Background()
	reader := &fakeReader{
		messages: []kafka.Message{
			{Value: []byte(`{"source": "RemoteOK", "source_id": 1, "company": "Acme", "role": "Backend Go Developer", "score": 1}`)},
		},
		err: errors.New("no more messages"),
	}
	jobStore := &fakeJobStore{
		createErr: errors.New("run create error"),
		listErr: nil,
		jobs: []domain.Job{},
	}
	consumer := NewConsumer(
		slog.New(slog.DiscardHandler),
		reader,
		jobStore,
	)

	err := consumer.Run(ctx)

	require.Error(t, err)
	assert.ErrorContains(t, err, "no more messages")
	assert.Equal(t, reader.calls, 1)
	require.Error(t, jobStore.createErr)
	require.Len(t, jobStore.jobs, 0)
}

func TestRun_FailJobConstruction(t *testing.T) {
	ctx := context.Background()
	reader := &fakeReader{
		messages: []kafka.Message{
			{Value: []byte(`{"source": "RemoteOK", "source_id": "1", "company": "", "role": "Backend Go Developer", "score": 1}`)},
		},
		err: errors.New("no more messages"),
	}
	jobStore := &fakeJobStore{
		createErr: errors.New("run job construction error"),
		listErr: nil,
		jobs: []domain.Job{},
	}
	consumer := NewConsumer(
		slog.New(slog.DiscardHandler),
		reader,
		jobStore,
	)

	err := consumer.Run(ctx)

	require.Error(t, err)
	assert.ErrorContains(t, err, "no more messages")
	assert.Equal(t, reader.calls, 1)
	require.Error(t, jobStore.createErr)
	require.Len(t, jobStore.jobs, 0)
}
