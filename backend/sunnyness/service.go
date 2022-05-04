package sunnyness

import (
	"time"

	"backend/interpolation"
	s "backend/shared"
	"backend/weatherapi"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const NumDecimalPlaces int = 1
const MinDegreeStep float32 = 0.1

//go:generate mockgen -destination=../mocks/mock_cache.go -package=mocks . Cache
type Cache interface {
	GetSunnyness(points *s.Point) (float32, error)
	SetSunnyness(points *s.Point) error
	SetSunnynesses(points []*s.Point) error
	CreateCompositeKey(point *s.Point) string
}

type SunnynessService interface {
	GetGrid(box s.Box, n s.NumPoints) (SunnynessGrid, error)
}

type Sunynessservice struct {
	cache                Cache
	api                  weatherapi.WeatherApi
	interpolationService interpolation.InterpolationService
	logger               log.Logger
}

func NewService(cache Cache, api weatherapi.WeatherApi, is interpolation.InterpolationService, logger log.Logger) SunnynessService {
	return &Sunynessservice{
		cache:                cache,
		api:                  api,
		interpolationService: is,
		logger:               logger,
	}
}

func (srv *Sunynessservice) GetGrid(b s.Box, n s.NumPoints) (SunnynessGrid, error) {
	start := time.Now()
	logger := log.With(srv.logger, "method", "GetGrid")
	level.Info(logger).Log("SERV", "getting grid")

	var diffLat float32 = b.BottomRightLat - b.TopLeftLat
	var diffLng float32 = b.BottomRightLng - b.TopLeftLng
	stepLat := s.Max(s.Abs(diffLat)/float32(n.Lat), MinDegreeStep)
	stepLng := s.Max(s.Abs(diffLng)/float32(n.Lng), MinDegreeStep)

	stepLat *= s.Sign(diffLat)
	stepLng *= s.Sign(diffLng)

	stepLat = s.FloorToDecimal(stepLat, NumDecimalPlaces)
	stepLng = s.FloorToDecimal(stepLng, NumDecimalPlaces)

	coords := CreateSnappedGridCoordinates(b, stepLat, stepLng)

	var queryPoints, cachePoints []*s.Point
	for _, c := range coords {
		sunnyness, err := srv.cache.GetSunnyness(c)
		if err != nil {
			level.Debug(logger).Log("SERV", "cache miss", "lat", c.Lat, "lng", c.Lng)
			queryPoints = append(queryPoints, c)
			continue
		}
		level.Debug(logger).Log("SERV", "cache hit", "lat", c.Lat, "lng", c.Lng)
		c.Val = sunnyness
		cachePoints = append(cachePoints, c)
	}

	srv.api.QueryPoints(queryPoints)
	go srv.cache.SetSunnynesses(queryPoints)

	queryPoints = append(queryPoints, cachePoints...)

	allPoints := srv.interpolationService.InterpolateGrid(queryPoints, b, n)

	elapsed := time.Since(start)
	level.Info(logger).Log("Elapsed time", elapsed)

	return SunnynessGrid{
		Points: allPoints,
	}, nil
}

func CreateSnappedGridCoordinates(b s.Box, stepLat float32, stepLng float32) []*s.Point {
	latStart := Snap(b.TopLeftLat, stepLat)
	lngStart := Snap(b.TopLeftLng, stepLng)
	var res []*s.Point
	for lat := latStart; s.Abs(lat) < s.Abs(b.BottomRightLat+stepLat); lat += stepLat {
		for lng := lngStart; s.Abs(lng) < s.Abs(b.BottomRightLng+stepLng); lng += stepLng {
			res = append(res, s.NewPoint(lat, lng))
		}
	}
	return res
}

// Snap snaps the given value to the
func Snap(val float32, step float32) float32 {
	flooredVal := s.FloorToDecimal(val, NumDecimalPlaces)
	var res float32
	if s.Sign(flooredVal) == s.Sign(step) {
		res = flooredVal - s.Mod(flooredVal, step)
	} else {
		res = flooredVal - (step + s.Mod(flooredVal, step))
	}
	return res
}

type SunnynessGrid struct {
	Points []*s.Point `json:"values"`
}
