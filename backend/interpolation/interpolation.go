package interpolation

import (
	s "backend/shared"

	"github.com/RobinRCM/sklearn/interpolate"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
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
		level.Debug(i.logger).Log("msg", "not interpolating result since queried points satisfy requested number of points", "len_queried_points", len(points), "requested_points", n.Lat*n.Lng)
		return points
	}

	var lats, lngs, vals []float64
	for _, p := range points {
		lats = append(lats, float64(p.Lat))
		lngs = append(lngs, float64(p.Lng))
		vals = append(vals, float64(p.Val))
	}
	ifunc := interpolate.Interp2d(lats, lngs, vals)

	res := []*s.Point{}
	// stepLng := (brlng - tllng) / float64(n.Lng)
	// stepLat := (brlat - tllat) / float64(n.Lat)
	stepLng, stepLat := s.CalculateStepSizes(b, n)

	var latStart, lngStart, latEnd, lngEnd float32
	latStart = s.Snap(s.Min(b.TopLeftLat, b.BottomRightLat), s.Abs(stepLat))
	latEnd = s.Snap(s.Max(b.TopLeftLat, b.BottomRightLat), -1*s.Abs(stepLat))

	lngStart = s.Snap(s.Min(b.TopLeftLng, b.BottomRightLng), s.Abs(stepLng))
	lngEnd = s.Snap(s.Max(b.TopLeftLng, b.BottomRightLng), -1*s.Abs(stepLng))

	var stepLatI float32 = (latEnd - latStart) / float32(n.Lat)
	var stepLngI float32 = (lngEnd - lngStart) / float32(n.Lng)

	for lat := latStart; lat <= latEnd; lat += s.Abs(stepLatI) {
		for lng := lngStart; lng <= lngEnd; lng += s.Abs(stepLngI) {
			var ip float32 = float32(ifunc(float64(lat), float64(lng)))
			ip = s.FloorToDecimal(float32(ip), 2)
			np := s.Point{
				Lat: s.FloorToDecimal(lat, s.Config.AppNumDecimalPlaces+2),
				Lng: s.FloorToDecimal(lng, s.Config.AppNumDecimalPlaces+2),
				Val: ip,
			}
			res = append(res, &np)
		}
	}

	level.Debug(i.logger).Log("msg", "interpolating result", "len_interpolated_points", len(res), "len_queried_points", len(points), "requested_points", n.Lat*n.Lng)
	return res
}
