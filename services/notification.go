package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/tbalthazar/onesignal-go"
	"log"
	"net/smtp"
)

func CreateNotificationByUserId(userID string, title string, message string, data models.RequestData) (res interface{}, err error) {
	notificationReq := &onesignal.NotificationRequest{
		AppID:     oneSignalAppID,
		Contents:  map[string]string{"en": message},
		Headings:  map[string]string{"en": title},
		IsAndroid: true,
		Data:      data,
		Tags: []interface{}{
			map[string]interface{}{
				"key":      "user_id",
				"relation": "=",
				"value":    userID,
			},
		},
	}

	createRes, _, err := oneSignalClient.Notifications.Create(notificationReq)
	if err != nil {
		return nil, err
	}

	return createRes, nil
}

func CreateEmailNotification(playerID string, emailSubject string, emailBody string) (err error) {
	from := "freeze.app.nextgen@gmail.com"
	password := "s3cr3tpassword"

	msg := "From: " + from + "\n" +
		"To: " + playerID + "\n" +
		"Subject: Your pin verification number\n\n" + emailBody

	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{playerID}, []byte(msg))

	if err != nil {
		log.Printf("smtp error : %s", err)
	}

	return err

}
