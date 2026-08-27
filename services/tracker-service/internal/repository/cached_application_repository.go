package repository

import (
	"context"
	"encoding/json"
	"errors"
	"jobradar/tracker-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

const applicationsCacheKey = "applications:all"

type CachedApplicationRepository struct {
    inner domain.ApplicationStore
    redis *redis.Client
    ttl   time.Duration
}

func NewCachedApplicationRepository(inner domain.ApplicationStore, redisClient *redis.Client, ttl time.Duration) *CachedApplicationRepository {
	return &CachedApplicationRepository {
		inner: inner,
		redis: redisClient,
		ttl: ttl,
	}
}

func (car *CachedApplicationRepository) List(ctx context.Context) ([]domain.Application, error) {
	redisCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	resp, err := car.redis.Get(redisCtx, applicationsCacheKey).Result()

	if err == nil {
		var applications []domain.Application
		err = json.Unmarshal([]byte(resp), &applications)

		if err != nil {
			// Log WARN
		} else{
			return applications, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		// Log WARN
	}

	applications, err := car.inner.List(ctx)

	if err != nil {
		return nil, err
	}

	jsonRecords, err := json.Marshal(&applications)

	if err != nil {
		// Log WARN
	}

	err = car.redis.Set(redisCtx, applicationsCacheKey, jsonRecords, car.ttl).Err()

	if err != nil {
		// Log WARN
	}

	return applications, nil
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
		// Log WARN
	}

	return nil
}
