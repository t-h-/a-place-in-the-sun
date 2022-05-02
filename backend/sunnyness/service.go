package sunnyness

import (
	"math"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const NumDecimalPlaces int = 1
const MinDegreeStep float32 = 0.1

//go:generate mockgen -destination=../mocks/mock_cache.go -package=mocks . Cache
type Cache interface {
	GetSunnyness(points *Point) (float32, error)
	SetSunnyness(points *Point) error
	SetSunnynesses(points []*Point) error
	CreateCompositeKey(point *Point) string
}

//go:generate mockgen -destination=../mocks/mock_api.go -package=mocks . WeatherApi
type WeatherApi interface {
	QueryPoints(points []*Point)
}

type SunnynessService interface {
	GetGrid(box Box, n NumPoints) (SunnynessGrid, error)
}

type Sunynessservice struct {
	cache  Cache
	api    WeatherApi
	logger log.Logger
}

func NewService(cache Cache, api WeatherApi, logger log.Logger) SunnynessService {
	return &Sunynessservice{
		cache:  cache,
		api:    api,
		logger: logger,
	}
}

func max(a float32, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func sign(f float32) float32 {
	if f == 0 {
		return 0
	}
	if f < 0 {
		return -1
	} else {
		return 1
	}
}

func abs(f float32) float32 {
	return f * sign(f)
}

func (s *Sunynessservice) GetGrid(b Box, n NumPoints) (SunnynessGrid, error) {
	start := time.Now()
	logger := log.With(s.logger, "method", "GetGrid")
	level.Info(logger).Log("SERV", "getting grid", "a", b.BottomRightLat)

	var diffLat float32 = b.BottomRightLat - b.TopLeftLat
	var diffLng float32 = b.BottomRightLng - b.TopLeftLng
	stepLat := max(abs(diffLat)/float32(n.Lat), MinDegreeStep)
	stepLng := max(abs(diffLng)/float32(n.Lng), MinDegreeStep)

	stepLat *= sign(diffLat)
	stepLng *= sign(diffLng)

	stepLat = floorToDecimal(stepLat, NumDecimalPlaces)
	stepLng = floorToDecimal(stepLng, NumDecimalPlaces)

	coords := CreateSnappedGridCoordinates(b, stepLat, stepLng)

	var queryPoints, cachePoints []*Point
	for _, c := range coords {
		sunnyness, err := s.cache.GetSunnyness(c)
		if err != nil {
			level.Debug(logger).Log("SERV", "cache miss", "lat", c.Lat, "lng", c.Lng)
			queryPoints = append(queryPoints, c)
			continue
		}
		level.Debug(logger).Log("SERV", "cache hit", "lat", c.Lat, "lng", c.Lng)
		c.Val = sunnyness
		cachePoints = append(cachePoints, c)
	}

	s.api.QueryPoints(queryPoints)
	s.cache.SetSunnynesses(queryPoints)

	elapsed := time.Since(start)
	// TODO interpolation service here, in case there are too few points
	level.Info(logger).Log("Elapsed time", elapsed)

	return SunnynessGrid{
		Points: append(cachePoints, queryPoints...),
	}, nil
}

func CreateSnappedGridCoordinates(b Box, stepLat float32, stepLng float32) []*Point {
	latStart := Snap(b.TopLeftLat, stepLat)
	lngStart := Snap(b.TopLeftLng, stepLng)
	var res []*Point
	for lat := latStart; abs(lat) < abs(b.BottomRightLat+stepLat); lat += stepLat {
		for lng := lngStart; abs(lng) < abs(b.BottomRightLng+stepLng); lng += stepLng {
			res = append(res, NewPoint(lat, lng))
		}
	}
	return res
}

func floorToDecimal(f float32, decimalPlaces int) float32 {
	factor := math.Pow(10, float64(decimalPlaces))
	return float32(math.Floor(float64(f*float32(factor))) / factor)
}

func mod(f1 float32, f2 float32) float32 {
	return float32(math.Mod(float64(f1), float64(f2)))
}

// Snap snaps the given value to the
func Snap(val float32, step float32) float32 {
	flooredVal := floorToDecimal(val, NumDecimalPlaces)
	var res float32
	if sign(flooredVal) == sign(step) {
		res = flooredVal - mod(flooredVal, step)
	} else {
		res = flooredVal - (step + mod(flooredVal, step))
	}
	return res
}

type Box struct {
	TopLeftLat     float32 `json:"top_left_lat"`
	TopLeftLng     float32 `json:"top_left_lng"`
	BottomRightLat float32 `json:"bottom_right_lat"`
	BottomRightLng float32 `json:"bottom_right_lng"`
}

type NumPoints struct {
	Lat int `json:"lat"`
	Lng int `json:"lng"`
}

type SunnynessGrid struct {
	Points []*Point `json:"values"`
}

type Point struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
	Val float32 `json:"val"`
}

func NewPoint(lat float32, lng float32) *Point {
	return &Point{
		Lat: floorToDecimal(lat, NumDecimalPlaces),
		Lng: floorToDecimal(lng, NumDecimalPlaces),
		Val: -1,
	}
}
