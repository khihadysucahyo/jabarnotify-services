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

//Notification schema
type Notification struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EmailAddress string             `json:"emailAddress,omitempty" bson:"emailAddress,omitempty"`
	PhoneNumber  string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Body         string             `json:"body,omitempty" bson:"body,omitempty" binding:"required"`
	Subject      string             `json:"subject,omitempty" bson:"subject,omitempty"`
	Type         string             `json:"type,omitempty" bson:"type,omitempty" binding:"required"`
	Status       string             `json:"status,omitempty" bson:"status,omitempty" binding:"required"`
	CreatedAt    time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty" binding:"required"`
	SendAt       time.Time          `json:"sendAt,omitempty" bson:"sendAt,omitempty"`
}

//SiteService describe the Stats service
type SiteService interface {
	GetNotification(ctx context.Context) ([]*Notification, error)
	CreateNotification(
		ctx context.Context,
		emailAddress string,
		phoneNumber string,
		body string,
		subject string,
		typ string) (*Notification, error)
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
func (s *basicService) GetNotification(ctx context.Context) ([]*Notification, error) {
	collection := s.DB.Collection("notifications")

	result, err := collection.Aggregate(ctx, mongo.Pipeline{})

	if err != nil {
		panic(err)
	}

	var data []*Notification

	if err = result.All(ctx, &data); err != nil {
		panic(err)
	}

	return data, nil
}

//CreateNotif display notif list
func (s *basicService) CreateNotification(
	ctx context.Context,
	emailAddress string,
	phoneNumber string,
	body string,
	subject string,
	typ string) (*Notification, error) {
	notification := &Notification{
		EmailAddress: emailAddress,
		PhoneNumber:  phoneNumber,
		Body:         body,
		Subject:      subject,
		Type:         typ,
		Status:       "sending",
	}

	collection := s.DB.Collection("notifications")
	insertResult, err := collection.InsertOne(context.TODO(), notification)

	if err != nil {
		return nil, err
	}

	fmt.Printf("type %T", insertResult)
	return notification, nil
}
