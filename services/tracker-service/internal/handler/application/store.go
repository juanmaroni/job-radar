package application

import (
	"context"
	"jobradar/tracker-service/internal/domain"
)

type ApplicationStore interface {
    Create(ctx context.Context, app *domain.Application) error
    List(ctx context.Context) ([]domain.Application, error)
}
