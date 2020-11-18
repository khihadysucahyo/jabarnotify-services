package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MessageNotification schema
type MessageNotification struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string             `json:"message,omitempty" bson:"message,omitempty" binding:"required"`
	Method    string             `json:"method,omitempty" bson:"method,omitempty" binding:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty" binding:"required"`
}

//SiteService describe the Stats service
type SiteService interface {
	GetMessageNotification(ctx context.Context) ([]*MessageNotification, error)
	CreateMessageNotification(ctx context.Context, msg string, mtd string) (*MessageNotification, error)
}

// NewSiteService returns a basic StatsService with all of the expected middlewares wired in.
func NewSiteService(logger log.Logger) SiteService {
	var svc SiteService
	svc = NewBasicService()
	svc = LoggingMiddleware(logger)(svc)
	return svc
}

// NewBasicService returns a naive, stateless implementation of StatsService.
func NewBasicService() SiteService {
	config, _ := utils.Initialize()
	return &basicService{
		DB: config.DB,
	}
}

type basicService struct {
	DB *mongo.Database
}

var (
	//ErrLoadNotif unable to find the requested team
	ErrLoadNotif = errors.New("error retriving notif")
)

//GetNotif display notif list
func (s *basicService) GetMessageNotification(ctx context.Context) ([]*MessageNotification, error) {
	collection := s.DB.Collection("messagenotifications")

	result, err := collection.Aggregate(ctx, mongo.Pipeline{})

	if err != nil {
		panic(err)
	}

	var data []*MessageNotification

	if err = result.All(ctx, &data); err != nil {
		panic(err)
	}

	return data, nil
}

//CreateNotif display notif list
func (s *basicService) CreateMessageNotification(ctx context.Context, msg string, mtd string) (*MessageNotification, error) {
	messageNotification := &MessageNotification{
		Message: msg,
		Method:  mtd,
	}

	collection := s.DB.Collection("messagenotifications")
	insertResult, err := collection.InsertOne(context.TODO(), messageNotification)

	if err != nil {
		return nil, err
	}

	fmt.Printf("type %T", insertResult)
	return messageNotification, nil
}
