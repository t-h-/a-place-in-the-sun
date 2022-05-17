package sunnyness

import (
	"backend/interpolation"
	s "backend/shared"
	"backend/weatherapi"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"
)

// TODO design error structure and handle them in transport layer
var ErrSthWentWrong = errors.New("something went wrong")

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

	stepLat, stepLng := s.CalculateStepSizes(b, n)

	coords := CreateSnappedGridCoordinates(b, stepLat, stepLng)

	var queryPoints, cachePoints []*s.Point
	for _, c := range coords {
		sunnyness, err := srv.cache.GetSunnyness(c)
		if err != nil {
			queryPoints = append(queryPoints, c)
			continue
		}
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
	var latStart, lngStart, latEnd, lngEnd float32

	latStart = s.Snap(s.Min(b.TopLeftLat, b.BottomRightLat), s.Abs(stepLat))
	latEnd = s.Snap(s.Max(b.TopLeftLat, b.BottomRightLat), -1*s.Abs(stepLat))

	lngStart = s.Snap(s.Min(b.TopLeftLng, b.BottomRightLng), s.Abs(stepLat))
	lngEnd = s.Snap(s.Max(b.TopLeftLng, b.BottomRightLng), -1*s.Abs(stepLat))

	// if stepLat >= 0 {
	// 	latStart = Snap(b.TopLeftLat, stepLat)
	// 	latEnd = Snap(b.BottomRightLat, -1*stepLat)
	// } else {
	// 	latStart = Snap(b.BottomRightLat, -1*stepLat)
	// 	latEnd = Snap(b.TopLeftLat, stepLat)
	// }

	// if stepLng >= 0 {
	// 	lngStart = Snap(b.TopLeftLng, stepLng)
	// 	lngEnd = Snap(b.BottomRightLng, -1*stepLng)
	// } else {
	// 	lngStart = Snap(b.BottomRightLng, -1*stepLng)
	// 	lngEnd = Snap(b.TopLeftLng, stepLng)
	// }

	var res []*s.Point
	for lat := latStart; lat <= latEnd; lat += s.Abs(stepLat) {
		for lng := lngStart; lng <= lngEnd; lng += s.Abs(stepLng) {
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

type SunnynessGrid struct {
	NumPoints int        `json:"numPoints"`
	Points    []*s.Point `json:"values"`
}
