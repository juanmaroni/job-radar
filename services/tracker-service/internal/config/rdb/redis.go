package rdb

import (
	"jobradar/tracker-service/internal/config/environment"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultTtl = 30

func CreateClient() (*redis.Client, error) {
	redisHost, err := environment.GetEnvVar("TRACKER_SERVICE_REDIS_HOST")

	if err != nil {
		return nil, err
	}

	redisPort, err := environment.GetEnvVar("TRACKER_SERVICE_REDIS_PORT")

	if err != nil {
		return nil, err
	}

	redisUri := redisHost + ":" + redisPort

	return redis.NewClient(
		&redis.Options{Addr: redisUri}), nil
}

func LoadTtl() (time.Duration, error) {
	rdbTtlStr, err := environment.GetEnvVar("TRACKER_SERVICE_REDIS_TTL")

	if err != nil {
		return time.Duration(defaultTtl) * time.Second, err
	}

	seconds, err := strconv.Atoi(rdbTtlStr)

	if err != nil {
		return time.Duration(defaultTtl) * time.Second, err
	} else {
		return time.Duration(seconds) * time.Second, nil
	}
}
