package sunnyness

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(svc SunnynessService, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	// TODO should be in main.go. need to see how to use this specific type of middlewear not on mux router but on http.Handler
	r.Use(commonMiddleware)

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	registerIncidentHandler := kithttp.NewServer(
		makeGetSunnynessGridEndpoint(svc),
		decodeGetSunnynessGridRequest,
		encodeGetSunnynessGridResponse,
		opts...,
	)

	r.Handle("/sunnyness/grid", registerIncidentHandler).Methods("POST")

	return r
}

func newServerFinalizer(logger kitlog.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		level.Info(logger).Log("status", code, "path", r.RequestURI, "method", r.Method, "params", fmt.Sprint(r.URL.Query()))
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

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrSthWentWrong:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

type errorer interface {
	error() error
} // TODO ?!?!?!?!

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
