package test

import (
	"fmt"
	"os"
	"testing"

	"backend/mocks"
	s "backend/shared"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
)

func TestSnap(t *testing.T) {
	var val float32 = 0.5
	var step float32 = 0.3
	snap := sunnyness.Snap(val, step)
	if snap != 0.3 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = -0.5
	step = 0.3
	snap = sunnyness.Snap(val, step)
	if snap != -0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = 0.5
	step = -0.3
	snap = sunnyness.Snap(val, step)
	if snap != 0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = -0.5
	step = -0.3
	snap = sunnyness.Snap(val, step)
	if snap != -0.3 {
		t.Fatalf(`Snap wrong %v`, snap)
	}
}

func TestCreateCoords(t *testing.T) {
	var flooredStepLat float32 = 0.5
	var flooredStepLng float32 = 0.5
	b := s.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	grid := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

func TestGetGrid(t *testing.T) {
	fmt.Println("TEST GetGrid()")

	svc, mock_api, mock_cache, mock_is := injectMocks(t)

	mock_is.EXPECT().InterpolateGrid(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]*s.Point{s.NewPoint(1, 1), s.NewPoint(2, 2)})
	mock_cache.EXPECT().GetSunnyness(gomock.Any()).AnyTimes().Return(float32(1.1), nil)
	mock_cache.EXPECT().SetSunnyness(gomock.Any()).AnyTimes().Return(nil)
	mock_cache.EXPECT().SetSunnynesses(gomock.Any()).AnyTimes().Return(nil)
	mock_api.EXPECT().QueryPoint(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return()

	b := s.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	n := s.NumPoints{Lat: 5, Lng: 5}
	grid, _ := svc.GetGrid(b, n)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

func injectMocks(t *testing.T) (sunnyness.SunnynessService, *mocks.MockWeatherService, *mocks.MockCache, *mocks.MockInterpolationService) {

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "debug", "true", "caller", log.DefaultCaller)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	api := mocks.NewMockWeatherService(ctrl)
	is := mocks.NewMockInterpolationService(ctrl)

	return sunnyness.NewService(cache, api, is, logger), api, cache, is
}
