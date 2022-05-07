package test

import (
	"os"
	"testing"

	s "backend/shared"
	"backend/sunnyness"
	"backend/weatherapi"

	"github.com/go-kit/log"
)

const ApiKey = "591b7934afcf484fa3191051223101"
const MaxRequestsPerSecond = 1
const MaxRequestBurst = 200

func TestQuery(t *testing.T) {
	var flooredStepLat float32 = 0.1
	var flooredStepLng float32 = 0.1
	b := s.Box{TopLeftLat: 9, TopLeftLng: 9, BottomRightLat: 100, BottomRightLng: 100}
	ps := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)

	api := createApi()

	api.QueryPoints(ps)
	// if len(ps) != 2 {
	// 	t.Fatalf(`res wrong %v`, ps)
	// }
	// fmt.Println(ps)
}

func createApi() weatherapi.WeatherApi {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	return weatherapi.NewApi(ApiKey, MaxRequestsPerSecond, MaxRequestBurst, logger)
}
