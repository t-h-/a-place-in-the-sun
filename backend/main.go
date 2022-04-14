package main

import (
	"net/http"
	"os"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "listen", "8083", "caller", log.DefaultCaller)

	cache := infra.NewInmemCache(logger)
	api := infra.NewApi(logger)

	svc := sunnyness.NewService(cache, api, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("err", http.ListenAndServe(":8083", router))
}
