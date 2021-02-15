package main

import "testing"

func TestMakeGetNotificationEndpoint(t *testing.T) {
	svc, _ := setup()
	makeGetNotificationEndpoint(svc)
}
func TestMakeCreateNotificationEndpoint(t *testing.T) {
	svc, _ := setup()
	makeCreateNotificationEndpoint(svc)
}

func TestMakeDetailNotificationEndpoint(t *testing.T) {
	svc, _ := setup()
	makeDetailNotificationEndpoint(svc)
}

func TestMakeGetNotificationSummaryEndpoint(t *testing.T) {
	svc, _ := setup()
	makeGetNotificationSummaryEndpoint(svc)
}

func TestMakeHealthCheckEndpoint(t *testing.T) {
	svc, _ := setup()
	makeHealthCheckEndpoint(svc)
}
