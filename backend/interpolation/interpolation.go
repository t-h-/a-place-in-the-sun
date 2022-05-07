package interpolation

import (
	s "backend/shared"

	"github.com/RobinRCM/sklearn/interpolate"
	"github.com/go-kit/kit/log"
)

//go:generate mockgen -destination=../mocks/mock_interpolation.go -package=mocks . InterpolationService
type InterpolationService interface {
	InterpolateGrid(points []*s.Point, b s.Box, n s.NumPoints) []*s.Point
}

type Interpolationservice struct {
	logger log.Logger
}

func NewInterpolationService(logger log.Logger) *Interpolationservice {
	return &Interpolationservice{
		logger: log.With(logger, "method", "bilinear interpolation"),
	}
}

// TODO consider moving Point to float64
func (i *Interpolationservice) InterpolateGrid(points []*s.Point, b s.Box, n s.NumPoints) []*s.Point {
	if len(points) > (n.Lat+2)*(n.Lng+2) {
		return points
	}

	brlat := float64(b.BottomRightLat)
	brlng := float64(b.BottomRightLng)
	tllat := float64(b.TopLeftLat)
	tllng := float64(b.TopLeftLng)

	var lats, lngs, vals []float64
	for _, p := range points {
		lats = append(lats, float64(p.Lat))
		lngs = append(lngs, float64(p.Lng))
		vals = append(vals, float64(p.Val))
	}
	ifunc := interpolate.Interp2d(lats, lngs, vals)

	res := []*s.Point{}
	stepLat := (brlat - tllat) / float64(n.Lat)
	stepLng := (brlng - tllng) / float64(n.Lng)
	for lat := tllat; lat <= brlat+stepLat; lat += stepLat {
		for lng := tllng; lng <= brlng+stepLng; lng += stepLng {
			ip := ifunc(lat, lng)
			ip = float64(s.FloorToDecimal(float32(ip), 2))
			np := s.Point{
				Lat: float32(lat),
				Lng: float32(lng),
				Val: float32(ip),
			}
			res = append(res, &np)
		}
	}

	return res
}
