package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//Endpoints holds all Stats Service enpoints
type Endpoints struct {
	GetNotification    endpoint.Endpoint
	CreateNotification endpoint.Endpoint
	DetailNotification endpoint.Endpoint
}

//MakeSiteEndpoints initialize all service Endpoints
func MakeSiteEndpoints(s SiteService) Endpoints {
	return Endpoints{
		GetNotification:    makeGetNotificationEndpoint(s),
		CreateNotification: makeCreateNotificationEndpoint(s),
		DetailNotification: makeDetailNotificationEndpoint(s),
	}
}

//NotificationRequest holds the request params for ListTables
type NotificationRequest struct {
	ID     string
	Method string
}

//NotificationReply holds the response params for ListTables
type NotificationReply struct {
	Items []map[string]interface{} `json:"items"`
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

func makeGetNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetNotification(ctx)
		return NotificationReply{Items: res, Err: err}, err
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
