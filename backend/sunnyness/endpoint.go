package sunnyness

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type GetSunnynessGridRequest struct {
	Box Box `json:"box"`
}

type GetSunnynessGridResponse struct {
	Grid SunnynessGrid `json:"grid,omitempty"`
	Err  string        `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

// makes function that decodes request to domain object and encodes domain response to request response
func makeGetSunnynessGridEndpoint(svc SunnynessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetSunnynessGridRequest)
		box := req.Box
		grid, err := svc.GetGrid(ctx, box)
		if err != nil {
			return GetSunnynessGridResponse{SunnynessGrid{}, err.Error()}, err
		}

		return GetSunnynessGridResponse{
			Grid: grid,
			Err:  "",
		}, err
	}
}
