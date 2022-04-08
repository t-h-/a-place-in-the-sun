package sunnyness

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type GetSunnynessGridRequest struct {
	TopLeftLat     float64 `json:"top_left_lat"`
	TopLeftLng     float64 `json:"top_left_lng"`
	BottomRightLat float64 `json:"bottom_right_lat"`
	BottomRightLng float64 `json:"bottom_right_lng"`
}

type GetSunnynessGridResponse struct {
	Values [][]int `json:"values,omitempty"`
	Err    string  `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

// makes function that decodes request to domain object and encodes domain response to request response
func makeGetSunnynessGridEndpoint(svc SunnynessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetSunnynessGridRequest)
		box := Box{
			TopLeftLat:     req.TopLeftLat,
			TopLeftLng:     req.TopLeftLng,
			BottomRightLat: req.BottomRightLat,
			BottomRightLng: req.BottomRightLng,
		}
		grid, err := svc.GetGrid(ctx, box)
		if err != nil {
			return GetSunnynessGridResponse{nil, err.Error()}, err
		}

		return GetSunnynessGridResponse{grid.Values, ""}, err
	}
}
