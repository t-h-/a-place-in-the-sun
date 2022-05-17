package test

import (
	"os"
	"testing"

	"backend/infra"
	s "backend/shared"

	"github.com/go-kit/log"
)

func TestCache(t *testing.T) {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	c, _ := infra.NewInmemCache(logger)

	c.SetSunnyness(s.NewPoint(1, 1))
	// ...
}
