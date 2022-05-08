package main

import (
	"net/http"
	"os"

	"backend/infra"
	"backend/interpolation"
	s "backend/shared"
	"backend/sunnyness"
	"backend/weatherapi"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	s.LoadConfigFromEnv()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	if s.Config.AppDebug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cache, _ := infra.NewInmemCache(logger)
	api := weatherapi.NewApi(logger)
	is := interpolation.NewInterpolationService(logger)

	svc := sunnyness.NewService(cache, api, is, logger)
	// svc := sunnyness.NewService(cache, api, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("msg", http.ListenAndServe(":8083", router))
}
