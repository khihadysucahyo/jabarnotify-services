package main

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func setupLoggerMiddleware() *log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	level.Info(logger).Log("msg", "service started v0.1")
	defer level.Info(logger).Log("msg", "service ended")

	return &logger
}

var logger = setupLoggerMiddleware()

func TestLoggerGetNotification(t *testing.T) {
	mw := LoggingMiddleware(*logger)(svc)
	mw.GetNotification(ctx, 1, 15)
}

func TestLoggerGetDetailNotif(t *testing.T) {
	mw := LoggingMiddleware(*logger)(svc)
	mw.DetailNotification(ctx, "602a2711236839792fb415ec")
}

func TestLoggerGetNotificationSummary(t *testing.T) {
	mw := LoggingMiddleware(*logger)(svc)
	mw.GetNotificationSummary(ctx)
}

func TestLoggerHealthCheck(t *testing.T) {
	mw := LoggingMiddleware(*logger)(svc)
	mw.HealthCheck(ctx)
}
