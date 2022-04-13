package main

import (
	"net/http"
	"os"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "listen", "8083", "caller", log.DefaultCaller)

	redisConn, err := connectToRedis()
	if err != nil {
		logger.Log("RDIS", "redis conn failed")
		//panic(err)
	}

	cache := infra.NewCache(redisConn, logger)
	api := infra.NewApi(logger)

	svc := sunnyness.NewService(cache, api, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("err", http.ListenAndServe(":8083", router))
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
