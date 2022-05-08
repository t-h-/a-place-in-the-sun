package test

import (
	"os"
	"testing"

	"backend/infra"

	"github.com/go-kit/log"
)

func TestReader(t *testing.T) {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	c, _ := infra.NewInmemCache(logger)

	n := float32(1.1)
	b, _ := c.Float32ToByte(float32(n))
	f, _ := c.ByteToFloat32(b)
	if f != n {
		t.Fatalf(`Float - Byte conversion wrong: %v - %v - %v`, n, b, f)
	}
}
