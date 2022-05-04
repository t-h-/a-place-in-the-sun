package shared

import (
	"math"
)

func Max(a float32, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func Sign(f float32) float32 {
	if f == 0 {
		return 0
	}
	if f < 0 {
		return -1
	} else {
		return 1
	}
}

func Abs(f float32) float32 {
	return f * Sign(f)
}

func FloorToDecimal(f float32, decimalPlaces int) float32 {
	factor := math.Pow(10, float64(decimalPlaces))
	return float32(math.Floor(float64(f*float32(factor))) / factor)
}

func Mod(f1 float32, f2 float32) float32 {
	return float32(math.Mod(float64(f1), float64(f2)))
}
