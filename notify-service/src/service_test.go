package main

import (
	"context"
	"testing"
)

func setup() (svc SiteService, ctx context.Context) {
	return NewBasicService(), context.Background()
}

var svc, ctx = setup()

func TestGetNotification(t *testing.T) {
	_, _, err := svc.GetNotification(ctx, 1, 15)

	if err != nil {
		t.Error("GetNotification service occured an error!")
	}
}

func TestDetailNotification(t *testing.T) {
	svc.DetailNotification(ctx, "602a2711236839792fb415ec")
}

func TestGetNotificationSummary(t *testing.T) {
	_, err := svc.GetNotificationSummary(ctx)

	if err != nil {
		t.Error("GetNotification service occured an error!")
	}
}

func TestHealthCheck(t *testing.T) {
	_, err := svc.HealthCheck(ctx)

	if err != nil {
		t.Error("GetNotification service occured an error!")
	}
}
