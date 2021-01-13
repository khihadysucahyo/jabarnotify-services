package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//Endpoints holds all Stats Service enpoints
type Endpoints struct {
	HealthCheck            endpoint.Endpoint
	GetNotification        endpoint.Endpoint
	GetNotificationSummary endpoint.Endpoint
	CreateNotification     endpoint.Endpoint
	DetailNotification     endpoint.Endpoint
}

//MakeSiteEndpoints initialize all service Endpoints
func MakeSiteEndpoints(s SiteService) Endpoints {
	return Endpoints{
		HealthCheck:            makeHealthCheckEndpoint(s),
		GetNotification:        makeGetNotificationEndpoint(s),
		GetNotificationSummary: makeGetNotificationSummaryEndpoint(s),
		CreateNotification:     makeCreateNotificationEndpoint(s),
		DetailNotification:     makeDetailNotificationEndpoint(s),
	}
}

//NotificationRequest holds the request params for ListTables
type NotificationRequest struct {
	ID      string
	Method  string
	Page    int
	PerPage int
}

//NotificationReply holds the response params for ListTables
type NotificationReply struct {
	Items []map[string]interface{} `json:"items"`
	Meta  *MetaData                `json:"meta"`
	Err   error                    `json:"err"`
}

//CreateNotificationRequest holds the request params for ListTables
type CreateNotificationRequest struct {
	Body       string
	Subject    string
	Type       string
	Recipients []*NotificationRecipient
}

//CreateNotificationReply holds the response params for ListTables
type CreateNotificationReply struct {
	Item *Notification `json:"item"`
	Err  error         `json:"err"`
}

//DetailNotificationReply holds the response params for ListTables
type DetailNotificationReply struct {
	Item map[string]interface{} `json:"item"`
	Err  error                  `json:"err"`
}

func makeHealthCheckEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.HealthCheck(ctx)
		return res, err
	}
}

func makeGetNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NotificationRequest)
		res, meta, err := s.GetNotification(ctx, req.Page, req.PerPage)
		return NotificationReply{Items: res, Meta: meta, Err: err}, err
	}
}

func makeCreateNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateNotificationRequest)
		result, err := s.CreateNotification(ctx, req.Body, req.Subject, req.Type, req.Recipients)
		return CreateNotificationReply{Item: result, Err: err}, err
	}
}

func makeDetailNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NotificationRequest)
		res, err := s.DetailNotification(ctx, req.ID)
		return DetailNotificationReply{Item: res, Err: err}, err
	}
}

func makeGetNotificationSummaryEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetNotificationSummary(ctx)
		return DetailNotificationReply{Item: res, Err: err}, err
	}
}
