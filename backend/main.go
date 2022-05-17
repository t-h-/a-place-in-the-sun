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

	mux := http.NewServeMux()
	mux.Handle("/sunnyness/", sunnyness.MakeHandler(svc, logger))
	http.Handle("/", accessControl(mux))

	logger.Log("msg", "HTTP", "addr", "8083")
	logger.Log("msg", http.ListenAndServe(":8083", nil))
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
