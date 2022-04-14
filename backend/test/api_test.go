package test

import (
	"fmt"
	"math"
	"os"
	"testing"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/log"
)

// TODO load test rate limiting
func TQuery(t *testing.T) {

	api := createApi()

	ps := []*sunnyness.Point{{Lat: 1.11, Lng: 2.22, Val: math.MaxFloat32}, {Lat: 3.33, Lng: 4.44, Val: math.MaxFloat32}}
	api.QueryPoints(ps)
	if len(ps) != 2 {
		t.Fatalf(`res wrong %v`, ps)
	}
	fmt.Println(ps)
}

func createApi() sunnyness.WeatherApi {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	return infra.NewApi(logger)
}
