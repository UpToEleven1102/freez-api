package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/tbalthazar/onesignal-go"
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

//func GetNotificationByUserId(userID string) (res interface{}, err error) {
//
//}
