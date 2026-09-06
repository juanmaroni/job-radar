package repository

import (
	"context"
	"errors"
	matchingproto "jobradar/proto/gen/go/matching/v1"
	"jobradar/tracker-service/internal/domain"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type fakeApplicationStore struct {
    createErr error
    listErr   error
    apps      []domain.Application
}

func (f *fakeApplicationStore) Create(ctx context.Context, app *domain.Application) error {
    if f.createErr != nil {
        return f.createErr
    }
	
    app.Id = int64(len(f.apps) + 1)
    f.apps = append(f.apps, *app)

    return nil
}

func (f *fakeApplicationStore) List(ctx context.Context) ([]domain.Application, error) {
    if f.listErr != nil {
        return nil, f.listErr
    }

    return f.apps, nil
}

type fakeMatchingClient struct {
    resp *matchingproto.MatchScoreResponse
    err  error
}

func (f *fakeMatchingClient) MatchScore(ctx context.Context, in *matchingproto.MatchScoreRequest, opts ...grpc.CallOption) (*matchingproto.MatchScoreResponse, error) {
    if f.err != nil {
        return nil, f.err
    }

    return f.resp, nil
}

func TestCreate_MatchingSuccess(t *testing.T) {
    var score int32 = 42
    fakeGrpc := &fakeMatchingClient{
        resp: &matchingproto.MatchScoreResponse{Score: score},
    }
    fakeInner := &fakeApplicationStore{}

    repo := NewMatchScoringApplicationRepository(
        slog.New(slog.DiscardHandler),
        fakeInner,
        fakeGrpc,
    )

    app, _ := domain.NewApplication("Acme", "Backend Engineer", "https://acme.com/jobs/1")

    err := repo.Create(context.Background(), app)

    require.NoError(t, err)
    require.NotNil(t, app.MatchScore)
    assert.Equal(t, score, *app.MatchScore)
	assert.Len(t, fakeInner.apps, 1)
}

func TestCreate_MatchingFail(t *testing.T) {
    fakeGrpc := &fakeMatchingClient{
        err: errors.New("matching service unavailable"),
    }
    fakeInner := &fakeApplicationStore{}

    repo := NewMatchScoringApplicationRepository(
        slog.New(slog.DiscardHandler),
        fakeInner,
        fakeGrpc,
    )

    app, _ := domain.NewApplication("Acme", "Backend Engineer", "https://acme.com/jobs/1")
    
	err := repo.Create(context.Background(), app)

    require.NoError(t, err)
    assert.Nil(t, app.MatchScore)
	assert.Len(t, fakeInner.apps, 1)
}

func TestList(t *testing.T) {
	now := time.Now()
	fakeInner := &fakeApplicationStore{
		apps: []domain.Application{
			{
				Id: 1,
				Company: "Acme",
				Role: "Backend Engineer",
				OfferUrl: "https://acme.com/jobs/1",
				Status: domain.StatusApplied,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				Id: 2,
				Company: "Google",
				Role: "Go Backend Engineer",
				OfferUrl: "https://acme.com/jobs/2",
				Status: domain.StatusInterview,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	repo := NewMatchScoringApplicationRepository(
        slog.New(slog.DiscardHandler),
        fakeInner,
        nil,
    )

	applications, err := repo.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, fakeInner.apps, 2)
	assert.Len(t, applications, 2)
    assert.Equal(t, fakeInner.apps, applications)
}
