package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//Endpoints holds all Stats Service enpoints
type Endpoints struct {
	GetNotification    endpoint.Endpoint
	CreateNotification endpoint.Endpoint
}

//MakeSiteEndpoints initialize all service Endpoints
func MakeSiteEndpoints(s SiteService) Endpoints {
	return Endpoints{
		GetNotification:    makeGetNotificationEndpoint(s),
		CreateNotification: makeCreateNotificationEndpoint(s),
	}
}

//NotificationRequest holds the request params for ListTables
type NotificationRequest struct {
	Method string
}

//NotificationReply holds the response params for ListTables
type NotificationReply struct {
	Items []*Notification `json:"items"`
	Err   error           `json:"err"`
}

//CreateNotificationRequest holds the request params for ListTables
type CreateNotificationRequest struct {
	EmailAddress string
	PhoneNumber  []string
	Body         string
	Subject      string
	Type         string
}

//CreateNotificationReply holds the response params for ListTables
type CreateNotificationReply struct {
	Item *Notification `json:"item"`
	Err  error         `json:"err"`
}

func makeGetNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetNotification(ctx)
		return NotificationReply{Items: res, Err: err}, nil
	}
}

func makeCreateNotificationEndpoint(s SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateNotificationRequest)
		result, err := s.CreateNotification(ctx, req.EmailAddress, req.PhoneNumber, req.Body, req.Subject, req.Type)
		return CreateNotificationReply{Item: result, Err: err}, nil
	}
}
