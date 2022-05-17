package test

import (
	"context"
	"os"
	"sync"
	"testing"

	s "backend/shared"
	"backend/weatherapi"

	"github.com/go-kit/log"
)

const ApiKey = "591b7934afcf484fa3191051223101"

func TestQueryPoint(t *testing.T) {
	srv := createWeatherService()

	cc := make(chan struct{}, 2)
	var wg sync.WaitGroup

	p := s.NewPoint(10, 10)

	wg.Add(1)
	cc <- struct{}{}
	srv.QueryPoint(p, &wg, cc)

	wg.Wait()
	// if len(ps) != 2 {
	// 	t.Fatalf(`res wrong %v`, ps)
	// }
	// fmt.Println(ps)
}

func createWeatherService() weatherapi.WeatherService {
	s.LoadConfigFromYaml("config.test.local.yml")
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	var srv weatherapi.WeatherService
	return weatherapi.NewProxyingMiddleware(ctx, logger)(srv)
}
