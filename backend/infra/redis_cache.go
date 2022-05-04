package infra

import (
	s "backend/shared"
	"backend/sunnyness"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-redis/redis"
)

const TTLMinutes = 60 * time.Minute

var (
	CacheErr      = errors.New("Unable to handle Repo Request")
	ErrIdNotFound = errors.New("Id not found")
)

type cache struct {
	client *redis.Client
	logger log.Logger
}

func NewCache(logger log.Logger) sunnyness.Cache {
	redisConn, err := connectToRedis()
	if err != nil {
		logger.Log("RDIS", "redis conn failed")
		//panic(err)
	}

	return &cache{
		client: redisConn,
		logger: log.With(logger, "cache", "redisTODO"),
	}
}

func (cache *cache) GetSunnyness(p *s.Point) (float32, error) {
	val, err := cache.client.Get(cache.CreateCompositeKey(p)).Float32()
	if err != nil {
		// TODO correct error handling
		return -1, ErrIdNotFound
	}

	return val, nil
}

func (cache *cache) SetSunnyness(p *s.Point) error {
	status_cmd := cache.client.Set(cache.CreateCompositeKey(p), p.Val, TTLMinutes)
	if err := status_cmd.Err(); err != nil {
		// TODO correct error handling
		return err
	}
	return nil
}

func (cache *cache) SetSunnynesses(points []*s.Point) error {
	for _, p := range points {
		err := cache.SetSunnyness(p)
		if err != nil {
			// TODO error handling
			return err
		}
	}
	return nil
}

func (cache *cache) CreateCompositeKey(p *s.Point) string {
	return fmt.Sprintf("%v:%v", p.Lat, p.Lng)
}

func connectToRedis() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := c.Ping().Result()

	return c, err
}
