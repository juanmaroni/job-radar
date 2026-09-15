package domain

import "context"

type JobStore interface {
	Create(ctx context.Context, job *Job) error
    List(ctx context.Context) ([]Job, error)
}
