package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// MakeHTTPHandler wires endpoints to the HTTP transport.
func MakeHTTPHandler(siteEndpoints Endpoints, logger log.Logger) http.Handler {
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "HEAD", "POST", "PUT"}),
	)

	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/notifications").Handler(kithttp.NewServer(
		siteEndpoints.GetNotification,
		decodeGetNotifRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST", "OPTIONS").Path("/notifications").Handler(kithttp.NewServer(
		siteEndpoints.CreateNotification,
		decodeCreateNotifRequest,
		encodeResponse,
		options...,
	))

	r.Use(cors)
	return r
}

func decodeGetNotifRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req NotificationRequest
	return req, nil
}

func decodeCreateNotifRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req CreateNotificationRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	fmt.Println("shshs")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrLoadNotif:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
