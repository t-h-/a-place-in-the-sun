package test

import (
	"fmt"
	"os"
	"testing"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

func TestSnap(t *testing.T) {
	var val float32 = 0.7
	var scale float32 = 0.3
	snap := sunnyness.Snap(val, scale)
	if snap != 0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}
}

func TestCreateCoords(t *testing.T) {
	var flooredStepLat float32 = 0.5
	var flooredStepLng float32 = 0.5
	b := sunnyness.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	grid := sunnyness.CreateSnappedGridCoordinates(b, flooredStepLat, flooredStepLng)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

func TestGetGrid(t *testing.T) {
	fmt.Println("TEST GetGrid()")

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "listen", "8083", "caller", log.DefaultCaller)

	redisConn, err := connectToRedis()
	if err != nil {
		logger.Log("RDIS", "redis conn failed")
		//panic(err)
	}

	cache := infra.NewCache(redisConn, logger)
	api := infra.NewApi(logger)

	srv := sunnyness.NewService(cache, api, logger)
	b := sunnyness.Box{TopLeftLat: 1.11, TopLeftLng: 1.11, BottomRightLat: 3.33, BottomRightLng: 3.33}
	n := sunnyness.NumPoints{Lat: 5, Lng: 5}
	grid, _ := srv.GetGrid(b, n)
	fmt.Println(grid)
	// if snap != 0.6 {
	// 	t.Fatalf(`Snap wrong %v`, snap)
	// }
}

// TERRIIIIBBLEEEE, in absence of go mocking knowledge...
func connectToRedis() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := c.Ping().Result()

	return c, err
}
