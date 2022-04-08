package sunnyness

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Cache interface {
	GetSunnyness(ctx context.Context, lat float64, lng float64) (int, error)             // TODO wrap into coord object
	SetSunnyness(ctx context.Context, lat float64, lng float64, val int) (string, error) // wrap into coord object
}

type SunnynessService interface {
	GetGrid(ctx context.Context, box Box) (SunnynessGrid, error)
}

type sunynessservice struct {
	cache  Cache
	logger log.Logger
}

func NewService(cache Cache, logger log.Logger) SunnynessService {
	return &sunynessservice{
		cache:  cache,
		logger: logger,
	}
}

func (s *sunynessservice) GetGrid(ctx context.Context, b Box) (SunnynessGrid, error) {
	logger := log.With(s.logger, "method", "Create") // ?!

	level.Info(logger).Log("info", nil) // TODO eh okay, stuff
	arr := [][]int{{int(b.TopLeftLat), int(b.TopLeftLng)}, {int(b.BottomRightLat), int(b.BottomRightLng)}}

	return SunnynessGrid{
		Values: arr,
	}, nil
}

type Box struct {
	TopLeftLat     float64 `json:"top_left_lat"`
	TopLeftLng     float64 `json:"top_left_lng"`
	BottomRightLat float64 `json:"bottom_right_lat"`
	BottomRightLng float64 `json:"bottom_right_lng"`
}

type SunnynessGrid struct {
	Values [][]int
}
