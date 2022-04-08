package sunnyness

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHttpServer(svc SunnynessService, logger kitlog.Logger) *mux.Router {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerFinalizer(newServerFinalizer(logger)),
	}
	getSunnynessGridHandler := kithttp.NewServer(
		makeGetSunnynessGridEndpoint(svc),
		decodeGetSunnynessGridRequest,
		encodeGetSunnynessGridResponse,
		options...,
	)
	r := mux.NewRouter()
	// r.Use(middleware.IsAuthenticatedMiddleware)
	r.Methods("GET").Path("/sunnyness/grid").
		Handler(getSunnynessGridHandler)
	return r
}

func newServerFinalizer(logger kitlog.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		logger.Log("status", code, "path", r.RequestURI, "method", r.Method, "params", fmt.Sprint(r.URL.Query()))
	}
}

func decodeGetSunnynessGridRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetSunnynessGridRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeGetSunnynessGridResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
