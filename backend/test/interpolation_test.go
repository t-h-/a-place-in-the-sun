package test

import (
	"backend/interpolation"
	s "backend/shared"
	"backend/sunnyness"
	"fmt"
	"os"
	"testing"

	"github.com/RobinRCM/sklearn/interpolate"
	"github.com/go-kit/log"
)

func TestInterpol(t *testing.T) {
	x := []float64{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4}
	y := []float64{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	z := []float64{0, 0, 0, 0, 0, 4, 6, 0, 0, 8, 10, 0, 0, 0, 0, 0}
	f := interpolate.Interp2d(x, y, z)
	ea := func(expected, actual float64) {
		if expected != actual {
			fmt.Printf("expected:%g actual:%g\n", expected, actual)
		}
	}
	a := f(2.5, 2.5)
	ea(7, a)
	ea(5.5, f(2.25, 2.25))

	if a == -1 {
		t.Fatalf(`Float - Byte conversion wrong:`)
	}
}

func TestInterpolSrv(t *testing.T) {
	var flooredStepLat float32 = 2
	var flooredStepLng float32 = 2
	b := s.Box{TopLeftLat: 10, TopLeftLng: 10, BottomRightLat: 20, BottomRightLng: 20}
	n := s.NumPoints{Lat: 10, Lng: 10}
	ps := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)

	for _, p := range ps {
		p.Val = (p.Lat + p.Lng) / 2
	}

	is := createInterpol()

	a := is.InterpolateGrid(ps, b, n)

	if a == nil {
		t.Fatalf(`Float - Byte conversion wrong:`)
	}
}

func createInterpol() interpolation.InterpolationService {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	return interpolation.NewService(logger)
}
