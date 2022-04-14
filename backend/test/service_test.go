package test

import (
	"fmt"
	"os"
	"testing"

	"backend/mocks"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
)

func TestSnap(t *testing.T) {
	var val float32 = 0.7
	var scale float32 = 0.3
	snap := sunnyness.Snap(val, scale)
	if snap != 0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}
}

func TestCreateCoords(t *testing.T) {
	var flooredStepLat float32 = 0.5
	var flooredStepLng float32 = 0.5
	b := sunnyness.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	grid := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

func TestGetGrid(t *testing.T) {
	fmt.Println("TEST GetGrid()")

	svc, mock_api, mock_cache := injectMocks(t)

	mock_cache.EXPECT().GetSunnyness(gomock.Any()).AnyTimes().Return(float32(1.1), nil)
	mock_cache.EXPECT().SetSunnyness(gomock.Any()).AnyTimes().Return("cool", nil)
	mock_cache.EXPECT().SetSunnynesses(gomock.Any()).AnyTimes().Return("cool", nil)
	mock_api.EXPECT().QueryPoints(gomock.Any()).AnyTimes().Return()

	b := sunnyness.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	n := sunnyness.NumPoints{Lat: 5, Lng: 5}
	grid, _ := svc.GetGrid(b, n)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

func injectMocks(t *testing.T) (sunnyness.SunnynessService, *mocks.MockWeatherApi, *mocks.MockCache) {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "listen", "8083", "caller", log.DefaultCaller)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	api := mocks.NewMockWeatherApi(ctrl)

	return sunnyness.NewService(cache, api, logger), api, cache
}
