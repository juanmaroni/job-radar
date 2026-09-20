package domain

import (
	"context"
)

type ApplicationStore interface {
    Create(ctx context.Context, app *Application) error
    List(ctx context.Context) ([]Application, error)
	UpdateStatus(ctx context.Context, id int64, status Status) error
}
