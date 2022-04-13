package test

import (
	"fmt"
	"os"
	"testing"

	"backend/infra"
	"backend/sunnyness"

	"github.com/go-kit/log"
)

func TestQuery(t *testing.T) {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	ps := []*sunnyness.Point{{Lat: 1.11, Lng: 2.22}, {Lat: 3.33, Lng: 4.44}}
	api := infra.NewApi(logger)
	res := api.QueryPoints(ps)
	if len(res) != 2 {
		t.Fatalf(`res wrong %v`, res)
	}
	fmt.Println(res)
}
