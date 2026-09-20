package repository

import (
	"context"
	"errors"
	"fmt"
	"jobradar/tracker-service/internal/domain"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

type ApplicationRepository struct {
	logger *slog.Logger
	pool *pgxpool.Pool
}

func NewApplicationRepository(logger *slog.Logger, pool *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository {
		logger: logger,
		pool: pool,
	}
}

func (ar *ApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	query := `
		INSERT INTO applications (company, role, offer_url, status, match_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	err := ar.pool.QueryRow(
		ctx, query, app.Company, app.Role, app.OfferUrl, app.Status, app.MatchScore, app.CreatedAt, app.UpdatedAt,
	).Scan(&app.Id)

	if err != nil {
		errInsert := fmt.Errorf("postgres insert application: %w", err)
		ar.logger.Error("postgres insert failed", slog.String("error", errInsert.Error()))
		
		return errInsert
	}

	return nil
}

func (ar *ApplicationRepository) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	query := `
		UPDATE applications 
		SET status = $1, updated_at = now()
		WHERE id = $3;
	`

	tag, err := ar.pool.Exec(ctx, query, status, id)

	if err != nil {
		errUpdate := fmt.Errorf("postgres update application id %d: %w", id, err)
		ar.logger.Error("postgres insert failed", slog.String("error", errUpdate.Error()))
		
		return errUpdate
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (ar *ApplicationRepository) List(ctx context.Context) ([]domain.Application, error) {
	query := `
		SELECT id, company, role, offer_url, status, match_score, created_at, updated_at
		FROM applications
		ORDER BY created_at DESC;
	`
	rows, err := ar.pool.Query(ctx, query)

	if err != nil {
		errSelect := fmt.Errorf("select application: %w", err)
		ar.logger.Error("postgres select failed", slog.String("error", errSelect.Error()))

		return nil, errSelect
	}

	lstApplications := []domain.Application{}

	for rows.Next() {
		var app domain.Application
		err := rows.Scan(&app.Id, &app.Company, &app.Role, &app.OfferUrl, &app.Status, &app.MatchScore, &app.CreatedAt, &app.UpdatedAt)

		if err != nil {
			ar.logger.Error("postgres scan failed", slog.String("error", err.Error()))

			return lstApplications, err
		}

		lstApplications = append(lstApplications, app)
	}

	if err := rows.Err(); err != nil {
		ar.logger.Error("postgres rows scan iteration failed", slog.String("error", err.Error()))

		return lstApplications, err
	}

	return lstApplications, nil
}
