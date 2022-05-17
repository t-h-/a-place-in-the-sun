package infra

import (
	s "backend/shared"
	"backend/sunnyness"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// TODO longterm: translate/unify all lower level errors to here defined errors
var (
	InmemCacheErr = errors.New("Unable to handle Inmem")
)

type inmemcache struct {
	sun    *bigcache.BigCache
	logger log.Logger
}

func NewInmemCache(logger log.Logger) (sunnyness.Cache, error) {
	lifeWindowSec := s.Config.CacheMaxLifeWindowSec
	bCache, err := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: time.Duration(lifeWindowSec) * time.Second,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: time.Duration(lifeWindowSec) * time.Second,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 4,

		// prints information about additional memory allocation
		Verbose: false,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 256,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	})
	if err != nil {
		level.Error(logger).Log("msg", "error creating cache", "error", err)
		return nil, err
	}

	return &inmemcache{
		sun:    bCache,
		logger: log.With(logger, "cache", "inmem"),
	}, nil
}

func (inmemcache *inmemcache) GetSunnyness(p *s.Point) (float32, error) {
	bs, err := inmemcache.sun.Get(inmemcache.CreateCompositeKey(p))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return -1, InmemCacheErr
		}

		return -1, err
	}

	v, err := s.ByteToFloat32(bs)
	if err != nil {
		level.Debug(inmemcache.logger).Log("msg", "byte to float conversion failed")
		return -1, err
	}
	level.Debug(inmemcache.logger).Log("msg", "returning value from cache", "lat", p.Lat, "lng", p.Lng, "val", p.Val)
	return v, nil
}

func (inmemcache *inmemcache) SetSunnyness(p *s.Point) error {
	if p.Val <= 0 {
		return nil
	}
	f, err := s.Float32ToByte(p.Val)
	if err != nil {
		level.Debug(inmemcache.logger).Log("msg", "float to byte conversion failed")
		return err
	}
	return inmemcache.sun.Set(inmemcache.CreateCompositeKey(p), f)
}

func (inmemcache *inmemcache) SetSunnynesses(points []*s.Point) error {
	for _, p := range points {
		err := inmemcache.SetSunnyness(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cache *inmemcache) CreateCompositeKey(p *s.Point) string {
	return fmt.Sprintf("%v:%v", p.Lat, p.Lng)
}
