package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var endpoints = MakeSiteEndpoints(svc)
var server = httptest.NewServer(MakeHTTPHandler(endpoints, *logger))

func TestGettingNotification(t *testing.T) {

	endpoints := MakeSiteEndpoints(svc)
	server := httptest.NewServer(MakeHTTPHandler(endpoints, *logger))
	defer server.Close()

	res := newHTTPServerCall(t, http.MethodGet, server.URL+"/notifications", nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK when getting list")
}

func TestCreateNotificationEndpoint(t *testing.T) {

	// Create Notof
	todo := map[string]interface{}{
		"subject": "this is subject",
		"body":    "this is request body",
		"type":    "whatsapp",
		"recipients": []*NotificationRecipient{
			&NotificationRecipient{
				Name:        "Test Name",
				PhoneNumber: "62888888888",
			},
		},
	}
	res := newHTTPServerCall(t, http.MethodPost, server.URL+"/notifications", todo)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK")

	var createNotificationReply CreateNotificationReply
	json.NewDecoder(res.Body).Decode(&createNotificationReply)
	// require.Equalf(t, todo., createNotificationReply.PhoneNumber, "todo set")
	// require.NotZerof(t, createNotificationReply.ID, "ID should be set")

	// var errorMap map[string]string
	// json.NewDecoder(res.Body).Decode(&errorMap)
	// require.EqualValuesf(t, "Not found", errorMap["error"], "Expected Not found error Todo")
}

func TestGetNotificationDetailEndpoint(t *testing.T) {

	endpoints := MakeSiteEndpoints(svc)
	server := httptest.NewServer(MakeHTTPHandler(endpoints, *logger))
	defer server.Close()

	res := newHTTPServerCall(t, http.MethodGet, server.URL+"/notifications/602a2711236839792fb415ec", nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK")
}

func TestGetNotificationSummaryEndpoint(t *testing.T) {

	endpoints := MakeSiteEndpoints(svc)
	server := httptest.NewServer(MakeHTTPHandler(endpoints, *logger))
	defer server.Close()

	res := newHTTPServerCall(t, http.MethodGet, server.URL+"/notifications/summary", nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK")
}

func TestHealthEndpoint(t *testing.T) {

	endpoints := MakeSiteEndpoints(svc)
	server := httptest.NewServer(MakeHTTPHandler(endpoints, *logger))
	defer server.Close()

	res := newHTTPServerCall(t, http.MethodGet, server.URL+"/health", nil)
	defer res.Body.Close()
	require.Equalf(t, http.StatusOK, res.StatusCode, "Expecting StatusOK")
}

// NewHTTPServerCall performs a http call
// It sets the request with all required headers. i.e. JWT token
func newHTTPServerCall(t *testing.T, httpMethod, url string, payload interface{}) *http.Response {
	var req *http.Request
	var err error

	if httpMethod == http.MethodGet || httpMethod == http.MethodDelete {
		req, err = http.NewRequest(httpMethod, url, nil)
	} else {
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(payload)
		req, err = http.NewRequest(httpMethod, url, b)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Authorization", newJWTToken(t))
	require.NoErrorf(t, err, "Error creating %s request", httpMethod)
	res, err := http.DefaultClient.Do(req)
	require.NoErrorf(t, err, "Error doing %s request to %s with payload %v", httpMethod, url, payload)
	return res
}
