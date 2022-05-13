package sunnyness

import (
	"backend/interpolation"
	s "backend/shared"
	"backend/weatherapi"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"
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
	weather              weatherapi.WeatherService
	interpolationService interpolation.InterpolationService
	logger               log.Logger
}

func NewService(cache Cache, api weatherapi.WeatherService, is interpolation.InterpolationService, logger log.Logger) SunnynessService {
	return &Sunynessservice{
		cache:                cache,
		weather:              api,
		interpolationService: is,
		logger:               logger,
	}
}

func (srv *Sunynessservice) GetGrid(b s.Box, n s.NumPoints) (SunnynessGrid, error) {
	start := time.Now()
	level.Info(srv.logger).Log("msg", "getting grid", "box", fmt.Sprintf("%v", b), "numPoints", fmt.Sprintf("%v", n))

	stepLat, stepLng := calculateStepSizes(b, n)

	coords := CreateSnappedGridCoordinates(b, stepLat, stepLng)

	var queryPoints, cachePoints []*s.Point
	for _, c := range coords {
		sunnyness, err := srv.cache.GetSunnyness(c)
		if err != nil {
			// level.Debug(srv.logger).Log("msg", "cache miss", "lat", c.Lat, "lng", c.Lng)
			queryPoints = append(queryPoints, c)
			continue
		}
		// level.Debug(srv.logger).Log("msg", "cache hit", "lat", c.Lat, "lng", c.Lng)
		c.Val = sunnyness
		cachePoints = append(cachePoints, c)
	}

	srv.QueryPoints(queryPoints)
	go srv.cache.SetSunnynesses(queryPoints)

	queryPoints = append(queryPoints, cachePoints...)

	allPoints := srv.interpolationService.InterpolateGrid(queryPoints, b, n)

	elapsed := time.Since(start)
	level.Debug(srv.logger).Log("msg", "Elapsed time", "total", elapsed) // TODO longterm: use instrumentation middleware for this

	return SunnynessGrid{
		NumPoints: len(allPoints),
		Points:    allPoints,
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

func (srv *Sunynessservice) QueryPoints(points []*s.Point) {
	var wg sync.WaitGroup
	rateLimiter := rate.NewLimiter(rate.Every(time.Duration(int64(1000/float32(s.Config.WeatherApiMaxRequestsPerSecond)))*time.Millisecond), s.Config.WeatherApiMaxRequestsPerSecond)
	var flag bool = true
	concurrencyTokens := make(chan struct{}, s.Config.WeatherApiMaxParallelRequests)
	for _, p := range points {
		ctx := context.Background()
		err := rateLimiter.Wait(ctx)
		if err != nil {
			level.Error(srv.logger).Log("msg", "error while waiting for ratelimiter", "error", err)
			return
		}
		concurrencyTokens <- struct{}{}
		wg.Add(1)
		go srv.weather.QueryPoint(p, &wg, concurrencyTokens)
	}
	go func() {
		wg.Wait()
		flag = false
		level.Debug(srv.logger).Log("msg", "done querying")
	}()

	// for debugging purposes...
	for flag {
		level.Debug(srv.logger).Log("msg", "open", "routines", runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond)
	}
}

// Snap snaps the given val to the closest multiple of step. If step is positive, then to the smaller multiple,
// if step is negative then to the bigger multiple. This is to make sure the starting point of our query is just a bit
// outside of the requested box, so that the frontend can display the heatmap neatly.
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

func calculateStepSizes(b s.Box, n s.NumPoints) (float32, float32) {
	var diffLat float32 = b.BottomRightLat - b.TopLeftLat
	var diffLng float32 = b.BottomRightLng - b.TopLeftLng
	stepLat := s.Max(s.Abs(diffLat)/float32(n.Lat), MinDegreeStep)
	stepLng := s.Max(s.Abs(diffLng)/float32(n.Lng), MinDegreeStep)

	stepLat *= s.Sign(diffLat)
	stepLng *= s.Sign(diffLng)

	stepLat = s.FloorToDecimal(stepLat, NumDecimalPlaces)
	stepLng = s.FloorToDecimal(stepLng, NumDecimalPlaces)

	return stepLat, stepLng
}

type SunnynessGrid struct {
	NumPoints int        `json:"numPoints"`
	Points    []*s.Point `json:"values"`
}
