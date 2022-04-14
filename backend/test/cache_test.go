package test

import (
	"testing"

	"backend/infra"
)

func TestReader(t *testing.T) {
	n := float32(1.1)
	b, _ := infra.Float32ToByte(float32(n))
	f, _ := infra.ByteToFloat32(b)
	if f != n {
		t.Fatalf(`Float - Byte conversion wrong: %v - %v - %v`, n, b, f)
	}
}
