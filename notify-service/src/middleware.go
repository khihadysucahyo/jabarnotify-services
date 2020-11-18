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

func (mw loggingMiddleware) GetMessageNotification(ctx context.Context) (t []*MessageNotification, err error) {
	defer func() {
		mw.logger.Log("method", "GetMessageNotification", "notif", "", "err", err)
	}()
	return mw.next.GetMessageNotification(ctx)
}

func (mw loggingMiddleware) CreateMessageNotification(ctx context.Context, msg string, mtd string) (t *MessageNotification, err error) {
	defer func() {
		mw.logger.Log("method", "CreateMessageNotification", "notif", "", "err", err)
	}()
	return mw.next.CreateMessageNotification(ctx, msg, mtd)
}
