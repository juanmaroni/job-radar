package domain

import (
	"context"
)

type JobSource interface {
    FetchJobs(ctx context.Context) ([]Job, error)
}
