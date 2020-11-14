package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/endpoint"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/service"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// MakeHTTPHandler wires endpoints to the HTTP transport.
func MakeHTTPHandler(siteEndpoints endpoint.Endpoints, logger log.Logger) http.Handler {

	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/message-notifications").Handler(kithttp.NewServer(
		siteEndpoints.GetMessageNotification,
		decodeGetNotifRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/message-notifications").Handler(kithttp.NewServer(
		siteEndpoints.CreateMessageNotification,
		decodeCreateNotifRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetNotifRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoint.MessageNotificationRequest
	return req, nil
}

func decodeCreateNotifRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoint.CreateMessageNotificationRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. Since we're using JSON, there's no reason to provide anything more specific.
// It's certainly possible to specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
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
	case service.ErrLoadNotif:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
