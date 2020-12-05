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

func (mw loggingMiddleware) GetNotification(ctx context.Context) (t []map[string]interface{}, err error) {
	defer func() {
		mw.logger.Log("method", "GetNotification", "notif", "", "err", err)
	}()
	return mw.next.GetNotification(ctx)
}

func (mw loggingMiddleware) DetailNotification(ctx context.Context, id string) (t map[string]interface{}, err error) {
	defer func() {
		mw.logger.Log("method", "GetNotification", "notif", "", "err", err)
	}()
	return mw.next.DetailNotification(ctx, id)
}

func (mw loggingMiddleware) CreateNotification(
	ctx context.Context,
	body string,
	subject string,
	typ string,
	recipients []*NotificationRecipient) (t *Notification, err error) {
	defer func() {
		mw.logger.Log("method", "CreateNotification", "notif", "", "err", err)
	}()
	return mw.next.CreateNotification(ctx, body, subject, typ, recipients)
}
