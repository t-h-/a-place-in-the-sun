package shared

import (
	"bytes"
	"encoding/binary"
	"math"
)

func Max(a float32, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func Min(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
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

// Snap snaps the given val to the closest multiple of step. If step is positive, then to the smaller multiple,
// if step is negative then to the bigger multiple. This is to make sure the starting point of our query is just a bit
// outside of the requested box, so that the frontend can display the heatmap neatly.
func Snap(val float32, step float32) float32 {
	flooredVal := FloorToDecimal(val, Config.AppNumDecimalPlaces)
	var res float32
	if Sign(flooredVal) == Sign(step) {
		res = flooredVal - Mod(flooredVal, step)
	} else {
		res = flooredVal - (step + Mod(flooredVal, step))
	}
	return res
}

func CalculateStepSizes(b Box, n NumPoints) (float32, float32) {
	var diffLat float32 = b.BottomRightLat - b.TopLeftLat
	var diffLng float32 = b.BottomRightLng - b.TopLeftLng
	stepLat := Max(Abs(diffLat)/float32(n.Lat), Config.AppMinDegreeStep)
	stepLng := Max(Abs(diffLng)/float32(n.Lng), Config.AppMinDegreeStep)

	stepLat *= Sign(diffLat)
	stepLng *= Sign(diffLng)

	stepLat = FloorToDecimal(stepLat, Config.AppNumDecimalPlaces)
	stepLng = FloorToDecimal(stepLng, Config.AppNumDecimalPlaces)

	return stepLat, stepLng
}

func Float32ToByte(f float32) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ByteToFloat32(bs []byte) (float32, error) {
	uintOfByte := binary.BigEndian.Uint32(bs)
	floatOfByte := math.Float32frombits(uintOfByte)

	return floatOfByte, nil
}
