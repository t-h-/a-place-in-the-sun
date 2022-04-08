package infra

import (
	"backend/sunnyness"
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-redis/redis"
)

var (
	CacheErr      = errors.New("Unable to handle Repo Request")
	ErrIdNotFound = errors.New("Id not found")
)

type cache struct {
	client *redis.Client
	logger log.Logger
}

func NewCache(cacheClient *redis.Client, logger log.Logger) (sunnyness.Cache, error) {
	return &cache{
		client: cacheClient,
		logger: log.With(logger, "cache", "redisTODO"),
	}, nil
}

func (cache *cache) GetSunnyness(ctx context.Context, lat float64, lng float64) (int, error) {
	val, err := cache.client.Get(fmt.Sprintf("%f", lat)).Int() // TODO add composite key
	if err != nil {
		return -1, ErrIdNotFound
	}

	return val, nil
}

func (cache *cache) SetSunnyness(ctx context.Context, lat float64, lng float64, val int) (string, error) {
	err := cache.client.Set(fmt.Sprintf("%f", lat), 100, 60) // TODO add composite key, make expiration time configurable
	if err != nil {
		// TODO error handling
	}

	return "sunnyness set succesfully", nil
}
