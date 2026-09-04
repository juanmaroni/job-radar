package repository

import (
	"context"
	matching "jobradar/proto/gen/go/matching/v1"
	"jobradar/tracker-service/internal/domain"
	"log/slog"
	"time"
)

type MatchScoringApplicationRepository struct {
	logger *slog.Logger
    inner domain.ApplicationStore
    grpc matching.MatchingServiceClient
}

func NewMatchScoringApplicationRepository(logger *slog.Logger, inner domain.ApplicationStore, grpcClient matching.MatchingServiceClient) *MatchScoringApplicationRepository {
    return &MatchScoringApplicationRepository{
        logger: logger,
        inner:  inner,
        grpc:   grpcClient,
    }
}

func (msar *MatchScoringApplicationRepository) List(ctx context.Context) ([]domain.Application, error) {
	return msar.inner.List(ctx)
}

func (msar *MatchScoringApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	grpcCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resp, err := msar.grpc.MatchScore(grpcCtx, &matching.MatchScoreRequest{Role: app.Role, Company: app.Company})

	if err != nil {
		msar.logger.Warn("grpc matchscore failed", slog.String("error", err.Error()))
	} else {
		app.MatchScore = &resp.Score
	}
	
	return msar.inner.Create(ctx, app)
}
