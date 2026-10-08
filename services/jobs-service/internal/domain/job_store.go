package domain

import "context"

type JobFilters struct {
	Days *int
}

type JobStore interface {
	Create(ctx context.Context, job *Job) error
    List(ctx context.Context, filters JobFilters) ([]Job, error)
}
