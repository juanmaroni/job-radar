package repository

import (
	"context"
	"fmt"
	"jobradar/tracker-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	pool *pgxpool.Pool
}

func NewApplicationRepository(pool *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository { pool: pool }
}

func (a *ApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	query := `
		INSERT INTO applications (company, role, offer_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	err := a.pool.QueryRow(
		ctx, query, app.Company, app.Role, app.OfferUrl, app.Status, app.CreatedAt, app.UpdatedAt,
	).Scan(&app.Id)

	if err != nil {
		return fmt.Errorf("insert application: %w", err)
	}

	return nil
}

func (a *ApplicationRepository)  List(ctx context.Context) ([]domain.Application, error) {
	query := `
		SELECT id, company, role, offer_url, status, created_at, updated_at
		FROM applications
		ORDER BY created_at DESC;
	`
	rows, err := a.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("select application: %w", err)
	}

	var lstApplications []domain.Application

	for rows.Next() {
		app := domain.Application{}
		err := rows.Scan(
			&app.Id, &app.Company, &app.Role, &app.OfferUrl, &app.Status, &app.CreatedAt, &app.UpdatedAt,
		)

		if err != nil {
			return lstApplications, err
		}

		lstApplications = append(lstApplications, app)
	}

	if err := rows.Err(); err != nil {
		return lstApplications, err
	}

	return lstApplications, nil
}
