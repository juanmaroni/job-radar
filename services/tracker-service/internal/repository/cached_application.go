package repository

import (
	"context"
	"encoding/json"
	"errors"
	"jobradar/tracker-service/internal/domain"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const applicationsCacheKey = "applications:all"

type CachedApplicationRepository struct {
	logger *slog.Logger
    inner domain.ApplicationStore
    redis *redis.Client
    ttl   time.Duration
}

func NewCachedApplicationRepository(logger *slog.Logger, inner domain.ApplicationStore, redisClient *redis.Client, ttl time.Duration) *CachedApplicationRepository {
	return &CachedApplicationRepository {
		logger: logger,
		inner: inner,
		redis: redisClient,
		ttl: ttl,
	}
}

func (car *CachedApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	redisCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := car.inner.Create(ctx, app)

	if err != nil {
		return err
	}

	err = car.redis.Del(redisCtx, applicationsCacheKey).Err()

	if err != nil {
		car.logger.Warn("redis del failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
	}

	return nil
}

func (car *CachedApplicationRepository) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	redisCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := car.inner.UpdateStatus(ctx, id, status)

	if err != nil {
		return err
	}

	err = car.redis.Del(redisCtx, applicationsCacheKey).Err()

	if err != nil {
		car.logger.Warn("redis del failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
	}

	return nil
}

func (car *CachedApplicationRepository) List(ctx context.Context) ([]domain.Application, error) {
	redisCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resp, err := car.redis.Get(redisCtx, applicationsCacheKey).Result()

	if err == nil {
		var applications []domain.Application
		err = json.Unmarshal([]byte(resp), &applications)

		if err != nil {
			car.logger.Warn("json unmarshall failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
		} else{
			return applications, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		car.logger.Warn("redis get failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
	}

	applications, err := car.inner.List(ctx)

	if err != nil {
		return nil, err
	}

	jsonRecords, err := json.Marshal(&applications)

	if err != nil {
		car.logger.Warn("json marshal failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
	}

	err = car.redis.Set(redisCtx, applicationsCacheKey, jsonRecords, car.ttl).Err()

	if err != nil {
		car.logger.Warn("redis set failed", slog.String("key", applicationsCacheKey), slog.String("error", err.Error()))
	}

	return applications, nil
}
