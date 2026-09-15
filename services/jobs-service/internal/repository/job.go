package repository

import (
	"context"
	"fmt"
	"jobradar/jobs-service/internal/domain"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
    logger *slog.Logger
    pool *pgxpool.Pool
}

func NewJobRepository(logger *slog.Logger, pool *pgxpool.Pool) *JobRepository {
    return &JobRepository{
		logger: logger,
		pool: pool,
	}
}

func (jr *JobRepository) Create(ctx context.Context, job *domain.Job) error {
    query := `
		INSERT INTO jobs (source, source_id, company, role, score, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (source, source_id) DO UPDATE
		SET company = EXCLUDED.company, role = EXCLUDED.role, score = EXCLUDED.score
		RETURNING id;
	`

	err := jr.pool.QueryRow(
		ctx, query, job.Source, job.SourceId, job.Company, job.Role, job.Score, job.CreatedAt,
	).Scan(&job.Id)

	if err != nil {
		errInsert := fmt.Errorf("postgres insert jobs: %w", err)
		jr.logger.Error("postgres insert failed", slog.String("error", errInsert.Error()))
		
		return errInsert
	}

	return nil
}

func (jr *JobRepository) List(ctx context.Context) ([]domain.Job, error) {
    query := `
		SELECT id, source, source_id, company, role, score, created_at
		FROM jobs
		ORDER BY created_at DESC;
	`
	rows, err := jr.pool.Query(ctx, query)

	if err != nil {
		errSelect := fmt.Errorf("select jobs: %w", err)
		jr.logger.Error("postgres select failed", slog.String("error", errSelect.Error()))

		return nil, errSelect
	}

	var lstJobs []domain.Job

	for rows.Next() {
		var job domain.Job
		err := rows.Scan(&job.Id, &job.Source, &job.SourceId, &job.Company, &job.Role, &job.Score, &job.CreatedAt)

		if err != nil {
			jr.logger.Error("postgres scan failed", slog.String("error", err.Error()))

			return lstJobs, err
		}

		lstJobs = append(lstJobs, job)
	}

	if err := rows.Err(); err != nil {
		jr.logger.Error("postgres rows scan iteration failed", slog.String("error", err.Error()))

		return lstJobs, err
	}

	return lstJobs, nil
}
