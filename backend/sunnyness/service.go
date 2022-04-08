package sunnyness

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Cache interface {
	GetSunnyness(ctx context.Context, lat float32, lng float32) (int, error)             // TODO wrap into coord object
	SetSunnyness(ctx context.Context, lat float32, lng float32, val int) (string, error) // wrap into coord object
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

	level.Info(logger).Log("SERV", "getting grid", "a", b.BottomRightLat) // TODO eh okay, stuff
	//arr := [][]int{{int(b.TopLeftLat), int(b.TopLeftLng)}, {int(b.BottomRightLat), int(b.BottomRightLng)}}
	arr := []Point{{Lat: b.TopLeftLat, Lng: b.TopLeftLng, Val: b.TopLeftLng}, {Lat: b.BottomRightLat, Lng: b.BottomRightLng, Val: b.BottomRightLng}}

	return SunnynessGrid{
		Values: arr,
	}, nil
}

type Box struct {
	TopLeftLat     float32 `json:"top_left_lat"`
	TopLeftLng     float32 `json:"top_left_lng"`
	BottomRightLat float32 `json:"bottom_right_lat"`
	BottomRightLng float32 `json:"bottom_right_lng"`
}

type SunnynessGrid struct {
	Values []Point `json:"values"`
}

type Point struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
	Val float32 `json:"val"`
}
