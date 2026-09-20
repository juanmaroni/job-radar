package repository

import (
	"context"
	matchingproto "jobradar/proto/gen/go/matching/v1"
	"jobradar/tracker-service/internal/domain"
	"log/slog"
	"time"
)

type MatchScoringApplicationRepository struct {
	logger *slog.Logger
    inner domain.ApplicationStore
    grpc matchingproto.MatchingServiceClient
}

func NewMatchScoringApplicationRepository(logger *slog.Logger, inner domain.ApplicationStore, grpcClient matchingproto.MatchingServiceClient) *MatchScoringApplicationRepository {
    return &MatchScoringApplicationRepository{
        logger: logger,
        inner:  inner,
        grpc:   grpcClient,
    }
}

func (msar *MatchScoringApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	if app.MatchScore == nil {
		grpcCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		resp, err := msar.grpc.MatchScore(grpcCtx, &matchingproto.MatchScoreRequest{Role: app.Role, Company: app.Company})

		if err != nil {
			msar.logger.Warn("grpc matchscore failed", slog.String("error", err.Error()))
		} else {
			app.MatchScore = &resp.Score
		}
	}
	
	return msar.inner.Create(ctx, app)
}

func (msar *MatchScoringApplicationRepository) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	return msar.inner.UpdateStatus(ctx, id, status)
}

func (msar *MatchScoringApplicationRepository) List(ctx context.Context) ([]domain.Application, error) {
	return msar.inner.List(ctx)
}
