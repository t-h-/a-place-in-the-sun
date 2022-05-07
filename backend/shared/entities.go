package shared

type Box struct {
	TopLeftLat     float32 `json:"top_left_lat"`
	TopLeftLng     float32 `json:"top_left_lng"`
	BottomRightLat float32 `json:"bottom_right_lat"`
	BottomRightLng float32 `json:"bottom_right_lng"`
}

type NumPoints struct {
	Lat int `json:"lat"`
	Lng int `json:"lng"`
}

type SunnynessGrid struct {
	Points []*Point `json:"values"`
}

type Point struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
	Val float32 `json:"val"`
}

func NewPoint(lat float32, lng float32) *Point {
	return &Point{
		Lat: FloorToDecimal(lat, 2), // TODO read this from global var
		Lng: FloorToDecimal(lng, 2),
		Val: -1,
	}
}
