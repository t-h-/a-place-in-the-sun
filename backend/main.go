package main

import (
	"net/http"
	"os"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// TODO get from .env / wrap in config package

const MaxRequestsPerSecond int = 2000
const MaxRequestBurst int = 2000
const CacheMaxLifeWindowSec int = 30
const debug bool = true

func main() {
	var ApiKeyy string = "591b7934afcf484fa3191051223101"

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	if debug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowDebug())
	}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cache := infra.NewInmemCache(CacheMaxLifeWindowSec, logger)
	api := infra.NewApi(ApiKeyy, MaxRequestsPerSecond, MaxRequestBurst, logger)

	svc := sunnyness.NewService(cache, api, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("12", http.ListenAndServe(":8083", router))
}
