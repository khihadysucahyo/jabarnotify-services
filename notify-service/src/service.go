package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userSession = make(map[string]map[string]interface{})

//Notification schema
type Notification struct {
	ID             primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Body           string                 `json:"body,omitempty" bson:"body,omitempty" binding:"required"`
	Subject        string                 `json:"subject,omitempty" bson:"subject,omitempty"`
	Type           string                 `json:"type,omitempty" bson:"type,omitempty" binding:"required"`
	RecipientTotal int                    `json:"recipientTotal,omitempty" bson:"recipientTotal,omitempty"`
	CreatedBy      map[string]interface{} `json:"createdBy,omitempty" bson:"createdBy,omitempty" binding:"required"`
	CreatedAt      time.Time              `json:"createdAt,omitempty" bson:"createdAt,omitempty" binding:"required"`
}

//NotificationRecipient schema
type NotificationRecipient struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NotificationID primitive.ObjectID `json:"body,notificationId" bson:"notificationId,omitempty" binding:"required"`
	Name           string             `json:"name,omitempty" bson:"name,omitempty"`
	EmailAddress   string             `json:"emailAddress,omitempty" bson:"emailAddress,omitempty"`
	PhoneNumber    string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`
	SendAt         time.Time          `json:"sendAt,omitempty" bson:"sendAt,omitempty"`
}

//SiteService describe the Stats service
type SiteService interface {
	GetNotification(ctx context.Context) ([]map[string]interface{}, error)
	CreateNotification(ctx context.Context,
		body string,
		subject string,
		typ string,
		recipients []*NotificationRecipient) (*Notification, error)
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
	//ErrRequest unable to find the requested team
	ErrRequest = errors.New("error request")
	//ErrUnauthorized unable to verify token
	ErrUnauthorized = errors.New("Unauthorized")
	//ErrExpiredToken handle expiredToken
	ErrExpiredToken = errors.New("Token is expired")
)

func getQueueName(typ string) string {
	queueName := ""

	switch typ {
	case "sms":
		queueName = "smsblast-queue"
	case "whatsapp":
		queueName = "wablast-queue"
	}

	return queueName
}

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

	svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"PhoneNumber": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(utils.CleanPhoneNumber(phoneNumber)),
			},
		},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(*confQueue.QueueUrl),
	})
}

//GetNotif display notif list
func (s *basicService) GetNotification(ctx context.Context) ([]map[string]interface{}, error) {
	collection := s.DB.Collection("notifications")
	sortStage := bson.D{{"$sort", bson.D{{"createdAt", -1}}}}

	result, err := collection.Aggregate(ctx, mongo.Pipeline{sortStage})
	fmt.Println(result)
	if err != nil {
		panic(err)
	}

	var data []map[string]interface{}

	if err = result.All(ctx, &data); err != nil {
		panic(err)
	}

	var showsLoaded []bson.M
	if err = result.All(ctx, &showsLoaded); err != nil {
		panic(err)
	}
	fmt.Println(showsLoaded)

	return data, nil
}

//CreateNotif display notif list
func (s *basicService) CreateNotification(
	ctx context.Context,
	body string,
	subject string,
	typ string,
	recipients []*NotificationRecipient) (*Notification, error) {
	notification := &Notification{
		Body:           body,
		Subject:        subject,
		Type:           typ,
		RecipientTotal: len(recipients),
		CreatedBy:      userSession["user"],
		CreatedAt:      time.Now(),
	}

	collection := s.DB.Collection("notifications")
	insertResult, err := collection.InsertOne(context.TODO(), notification)

	if err != nil {
		return nil, err
	}

	queueName := getQueueName(typ)

	if queueName == "" || len(recipients) < 1 {
		return nil, ErrRequest
	}

	for _, recipient := range recipients {

		notificationRecipient := &NotificationRecipient{
			NotificationID: insertResult.InsertedID.(primitive.ObjectID),
			Name:           recipient.Name,
			EmailAddress:   recipient.EmailAddress,
			PhoneNumber:    recipient.PhoneNumber,
			Status:         "sent",
		}

		collection := s.DB.Collection("notificationrecipients")
		collection.InsertOne(context.TODO(), notificationRecipient)

		messg := body
		messg = strings.ReplaceAll(messg, "{NAME}", recipient.Name)
		messg = strings.ReplaceAll(messg, "{PHONE_NUMBER}", recipient.PhoneNumber)

		pushNotifToPhoneNumber(queueName, recipient.PhoneNumber, messg)
	}

	fmt.Printf("type %T", insertResult)
	return notification, nil
}
