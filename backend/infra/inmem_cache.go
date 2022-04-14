package infra

import (
	"backend/sunnyness"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-kit/log"
)

var (
	InmemCacheErr = errors.New("Unable to handle Inmem")
)

type inmemcache struct {
	sun    *bigcache.BigCache
	logger log.Logger
}

func NewInmemCache(logger log.Logger) *inmemcache {
	bCache, err := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 30 * time.Second,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 30 * time.Second,

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
		// return nil, fmt.Errorf("new big cache: %w", err)
		// TODO error handling, panic?
	}

	return &inmemcache{
		sun:    bCache,
		logger: log.With(logger, "cache", "inmem"), // TODO correct logging (all of the project)
	}
}

func (inmemcache *inmemcache) GetSunnyness(p *sunnyness.Point) (float32, error) {
	bs, err := inmemcache.sun.Get(inmemcache.CreateCompositeKey(p))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return -1, InmemCacheErr
		}

		return -1, fmt.Errorf("get: %w", err)
	}

	return ByteToFloat32(bs)
}

func (inmemcache *inmemcache) SetSunnyness(p *sunnyness.Point) error {
	f, err := Float32ToByte(p.Val)
	if err != nil {
		// TODO error handling
	}
	return inmemcache.sun.Set(inmemcache.CreateCompositeKey(p), f)
}

func (inmemcache *inmemcache) SetSunnynesses(points []*sunnyness.Point) error {
	for _, p := range points {
		err := inmemcache.SetSunnyness(p)
		if err != nil {
			// TODO error handling
			return err
		}
	}

	return nil
}

func (cache *inmemcache) CreateCompositeKey(p *sunnyness.Point) string {
	return fmt.Sprintf("%v:%v", p.Lat, p.Lng)
}

func Float32ToByte(f float32) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		// TODO error handling
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToFloat32(bs []byte) (float32, error) {
	fdsa := binary.BigEndian.Uint32(bs)
	asdf := math.Float32frombits(fdsa)

	// var f float32
	// var r = bytes.NewReader(bs)
	// err := binary.Read(r, binary.BigEndian, f)
	// if err != nil {
	// 	// TODO error handling
	// 	fmt.Println("binary.Read failed:", err)
	// 	return -1, err
	// }
	return asdf, nil
}
