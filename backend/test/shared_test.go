package test

import (
	s "backend/shared"
	"testing"
)

func TestReader(t *testing.T) {

	n := float32(1.1)
	b, _ := s.Float32ToByte(float32(n))
	f, _ := s.ByteToFloat32(b)
	if f != n {
		t.Fatalf(`Float - Byte conversion wrong: %v - %v - %v`, n, b, f)
	}
}

func TestSnap(t *testing.T) {
	var val float32 = 0.5
	var step float32 = 0.3
	snap := s.Snap(val, step)
	if snap != 0.3 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = -0.5
	step = 0.3
	snap = s.Snap(val, step)
	if snap != -0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = 0.5
	step = -0.3
	snap = s.Snap(val, step)
	if snap != 0.6 {
		t.Fatalf(`Snap wrong %v`, snap)
	}

	val = -0.5
	step = -0.3
	snap = s.Snap(val, step)
	if snap != -0.3 {
		t.Fatalf(`Snap wrong %v`, snap)
	}
}
