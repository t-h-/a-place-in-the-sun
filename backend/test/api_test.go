package test

import (
	"context"
	"os"
	"sync"
	"testing"

	s "backend/shared"
	"backend/sunnyness"
	"backend/weatherapi"

	"github.com/go-kit/log"
)

const ApiKey = "591b7934afcf484fa3191051223101"

func TestQuery(t *testing.T) {
	var flooredStepLat float32 = 0.1
	var flooredStepLng float32 = 0.1
	b := s.Box{TopLeftLat: 9, TopLeftLng: 9, BottomRightLat: 10, BottomRightLng: 10}
	ps := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)

	srv := createWeatherService()

	cc := make(chan struct{}, 2)
	var wg sync.WaitGroup

	srv.QueryPoint(ps[0], &wg, cc)
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
