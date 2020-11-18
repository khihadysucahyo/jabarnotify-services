package main

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(SiteService) SiteService

// LoggingMiddleware takes a logger as a dependency and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next SiteService) SiteService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   SiteService
}

func (mw loggingMiddleware) GetNotification(ctx context.Context) (t []*Notification, err error) {
	defer func() {
		mw.logger.Log("method", "GetNotification", "notif", "", "err", err)
	}()
	return mw.next.GetNotification(ctx)
}

func (mw loggingMiddleware) CreateNotification(
	ctx context.Context,
	emailAddress string,
	phoneNumber string,
	body string,
	subject string,
	typ string) (t *Notification, err error) {
	defer func() {
		mw.logger.Log("method", "CreateNotification", "notif", "", "err", err)
	}()
	return mw.next.CreateNotification(ctx, emailAddress, phoneNumber, body, subject, typ)
}