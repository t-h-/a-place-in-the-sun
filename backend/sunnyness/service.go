package sunnyness

import (
	"math"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const NumDecimalPlaces int = 1
const MinDegreeStep float64 = 0.1

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

func (s *Sunynessservice) GetGrid(b Box, n NumPoints) (SunnynessGrid, error) {
	logger := log.With(s.logger, "method", "Create") // ?!
	level.Info(logger).Log("SERV", "getting grid", "a", b.BottomRightLat)

	stepLat := float32(math.Max(float64((b.TopLeftLat-b.BottomRightLat)/float32(n.Lat)), MinDegreeStep)) // write own max method for float32
	stepLng := float32(math.Max(float64((b.TopLeftLng-b.BottomRightLng)/float32(n.Lng)), MinDegreeStep))

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

	// TODO interpolation service here, in case there are too few points

	return SunnynessGrid{
		Points: append(cachePoints, queryPoints...),
	}, nil
}

func CreateSnappedGridCoordinates(b Box, stepLat float32, stepLng float32) []*Point {
	latStart := Snap(floorToDecimal(b.TopLeftLat, NumDecimalPlaces), stepLat)
	lngStart := Snap(floorToDecimal(b.TopLeftLng, NumDecimalPlaces), stepLng)
	var res []*Point
	for lat := latStart; lat < b.BottomRightLat+stepLat; lat += stepLat {
		for lng := lngStart; lng < b.BottomRightLng; lng += stepLng {
			res = append(res, NewPoint(lat, lng))
		}
	}
	return res
}

func floorToDecimal(f float32, decimalPlaces int) float32 {
	factor := math.Pow(10, float64(decimalPlaces))
	return float32(math.Floor(float64(f*float32(factor))) / factor)
}

func Snap(val float32, scale float32) float32 {
	flooredVal := floorToDecimal(val, NumDecimalPlaces)
	res := flooredVal - float32(math.Mod(float64(flooredVal), float64(scale)))
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
		Val: float32(math.NaN()),
	}
}
