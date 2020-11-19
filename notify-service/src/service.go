package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func pushNotifToPhoneNumber(queueName string, phoneNumber string, body string) {
	sess := session.New(&aws.Config{
		Region: aws.String(utils.GetEnv("AWS_DEFAULT_REGION")),
		Credentials: credentials.NewStaticCredentials(
			utils.GetEnv("AWS_ACCESS_KEY_ID"),
			utils.GetEnv("AWS_SECRET_ACCESS_KEY"), "",
		),
		MaxRetries: aws.Int(5),
	})

	svc := sqs.New(sess)

	confQueue, _ := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"PhoneNumber": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(phoneNumber),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(*confQueue.QueueUrl),
	})

	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println("Success", *result.MessageId)

}

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

	if typ == "sms" {
		pushNotifToPhoneNumber("smsblast-queue", phoneNumber, body)
	}

	fmt.Printf("type %T", insertResult)
	return notification, nil
}
