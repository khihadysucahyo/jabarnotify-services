package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/pkg/service"
)

//Endpoints holds all Stats Service enpoints
type Endpoints struct {
	GetMessageNotification    endpoint.Endpoint
	CreateMessageNotification endpoint.Endpoint
}

//MakeSiteEndpoints initialize all service Endpoints
func MakeSiteEndpoints(s service.SiteService) Endpoints {
	return Endpoints{
		GetMessageNotification:    makeGetMessageNotificationEndpoint(s),
		CreateMessageNotification: makeCreateMessageNotificationEndpoint(s),
	}
}

//MessageNotificationRequest holds the request params for ListTables
type MessageNotificationRequest struct {
	Method string
}

//MessageNotificationReply holds the response params for ListTables
type MessageNotificationReply struct {
	Items []*service.MessageNotification `json:"items"`
	Err   error                          `json:"err"`
}

//CreateMessageNotificationRequest holds the request params for ListTables
type CreateMessageNotificationRequest struct {
	Message string `json:"message"`
	Method  string `json:"method"`
}

//CreateMessageNotificationReply holds the response params for ListTables
type CreateMessageNotificationReply struct {
	Item *service.MessageNotification `json:"item"`
	Err  error                        `json:"err"`
}

func makeGetMessageNotificationEndpoint(s service.SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetMessageNotification(ctx)
		return MessageNotificationReply{Items: res, Err: err}, nil
	}
}

func makeCreateMessageNotificationEndpoint(s service.SiteService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateMessageNotificationRequest)
		result, err := s.CreateMessageNotification(ctx, req.Message, req.Method)
		return CreateMessageNotificationReply{Item: result, Err: err}, nil
	}
}
