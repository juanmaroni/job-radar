package domain

import "context"

type JobPublisher interface {
    Publish(ctx context.Context, job Job) error
}
