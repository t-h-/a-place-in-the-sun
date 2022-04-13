package infra

import (
	"backend/sunnyness"
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

func NewCache(cacheClient *redis.Client, logger log.Logger) sunnyness.Cache {
	return &cache{
		client: cacheClient,
		logger: log.With(logger, "cache", "redisTODO"),
	}
}

func (cache *cache) GetSunnyness(p *sunnyness.Point) (float32, error) {
	val, err := cache.client.Get(fmt.Sprintf("%f", p.Lat)).Float32() // TODO add composite key
	if err != nil {
		// TODO correct error handling
		return -1, ErrIdNotFound
	}

	return val, nil
}

func (cache *cache) SetSunnyness(p *sunnyness.Point) (string, error) {
	err := cache.client.Set(fmt.Sprintf("%f", p.Lat), 100, 60) // TODO add composite key, make expiration time configurable
	if err != nil {
		// TODO error handling
	}

	return "sunnyness set succesfully", nil
}

func (cache *cache) SetSunnynesses(points []*sunnyness.Point) (string, error) {
	for _, p := range points {
		_, err := cache.SetSunnyness(p)
		if err != nil {
			// TODO error handling
		}
	}

	return "sunnyness set succesfully", nil
}
