package main

import (
	"strings"
	"testing"

	"github.com/jabardigitalservice/jabarnotify-services/notify-service/src/utils"
)

func TestCreateNotification(t *testing.T) {
	// payload
	typ := "whatsapp"
	body := "Running Go test: send whatsapp passed"
	recipients := []*NotificationRecipient{
		&NotificationRecipient{
			Name:        "GO TEST",
			PhoneNumber: utils.GetEnv("PHONE_NUMBER_TESTER"),
		},
	}

	queueName := getQueueName(typ)

	for _, recipient := range recipients {
		messg := body
		messg = strings.ReplaceAll(messg, "{NAME}", recipient.Name)
		messg = strings.ReplaceAll(messg, "{PHONE_NUMBER}", recipient.PhoneNumber)

		// pushNotifToPhoneNumber(queueName, recipient.PhoneNumber, messg)
	}
}
