package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jabardigitalservice/jabarnotify-services/notify-service/src/utils"
	"github.com/stretchr/testify/require"
)

func TestSendWhatsapp(t *testing.T) {
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
		QueueName: aws.String("wablast-queue"),
	})

	res, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"PhoneNumber": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(utils.GetEnv("PHONE_NUMBER_TESTER")),
			},
		},
		MessageBody: aws.String("running go test: send whatsapp passed"),
		QueueUrl:    aws.String(*confQueue.QueueUrl),
	})

	require.NoError(t, err)
	require.NotNil(t, res)
}
