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

	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	cache, err := infra.NewCache(redisConn, logger)
	if err != nil {
		logger.Log("TODO SOMETHING!")
	}

	svc := sunnyness.NewService(cache, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("err", http.ListenAndServe(":8083", router))
}
