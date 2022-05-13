package main

import (
	"backend/infra"
	"backend/interpolation"
	s "backend/shared"
	"backend/sunnyness"
	"backend/weatherapi"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	s.LoadConfigFromEnv()

	var logger log.Logger
	var ctx = context.Background()
	logger = log.NewLogfmtLogger(os.Stderr)
	if s.Config.AppDebug {
		logger = level.NewFilter(logger, level.AllowDebug())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	level.Info(logger).Log("config", fmt.Sprintf("%v", s.Config))

	cache, _ := infra.NewInmemCache(logger)

	var weatherApi weatherapi.WeatherService
	weatherApi = weatherapi.NewProxyingMiddleware(ctx, logger)(weatherApi)

	is := interpolation.NewInterpolationService(logger)

	svc := sunnyness.NewService(cache, weatherApi, is, logger)
	// svc := sunnyness.NewService(cache, api, logger)
	router := sunnyness.NewHttpServer(svc, logger)
	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("msg", http.ListenAndServe(":8083", router))
}
