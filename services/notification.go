package services

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/tbalthazar/onesignal-go"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)


func InsertMerchantNotification(merchantID string, activityType int, sourceID int64, message string) error {
	_, err := DB.Exec(`INSERT INTO merchant_notification (merchant_id, activity_type, source_id, message) VALUES (?, ?, ?, ?)`, merchantID, activityType, sourceID, message)
	if err != nil {
		log.Println(err)
	}
	return err
}

func InsertUserNotification(userID string, activityType int, sourceID int64, message string) error {
	_, err := DB.Exec(`INSERT INTO user_notification (user_id, activity_type, source_id, message) VALUES (?, ?, ?, ?)`, userID, activityType, sourceID, message)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateUserNotification(notification models.UserNotification) error {
	_, err := DB.Exec(`UPDATE user_notification SET unread=? WHERE id=?`, notification.UnRead, notification.ID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func UpdateMerchantNotification(notification models.MerchantNotification) error {
	_, err := DB.Exec(`UPDATE merchant_notification SET unread=? WHERE id=?`, notification.UnRead, notification.ID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func GetUserNotificationById(id int64) (interface{}, error){
	r, err := DB.Query(`SELECT u.id, ts, user_id, type, source_id, unread, message 
								FROM user_notification u
								LEFT JOIN activity_type a on u.activity_type = a.id
								WHERE u.id=?`, id)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer r.Close()
	var notification models.UserNotification
	if r.Next()  {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.UserID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		notificationInfo := models.UserNotificationInfo{ID: notification.ID, TimeStamp:notification.TimeStamp, UserID:notification.UserID, ActivityType:notification.ActivityType, UnRead:notification.UnRead, Message:notification.Message}
		switch notification.ActivityType {
		case "request":
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		}
		return notificationInfo, nil
	}
	return nil, nil
}

func GetMerchantNotificationById(id int64) (interface{}, error) {
	r, err := DB.Query(`SELECT m.id, ts, merchant_id, type, source_id, unread, message 
								FROM merchant_notification m
								LEFT JOIN activity_type a on m.activity_type = a.id
								WHERE m.id=?`, id)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer r.Close()

	var notification models.MerchantNotification
	if r.Next()  {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.MerchantID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		notificationInfo := models.MerchantNotificationInfo{ID: notification.ID, TimeStamp:notification.TimeStamp, MerchantID:notification.MerchantID, ActivityType:notification.ActivityType, UnRead:notification.UnRead, Message:notification.Message}
		switch notification.ActivityType {
		case "request":
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		}

		return notificationInfo, nil
	}
	return nil, nil
}

func GetMerchantNotifications(merchantID string) (notifications []interface{}, err error) {
	notifications = []interface{}{}
	r, err := DB.Query(`SELECT m.id, ts, merchant_id, type, source_id, unread, message 
								FROM merchant_notification m
								LEFT JOIN activity_type a on m.activity_type = a.id
								WHERE merchant_id=?`, merchantID)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer r.Close()
	var notification models.MerchantNotification
	for r.Next()  {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.MerchantID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
		}

		notificationInfo := models.MerchantNotificationInfo{ID: notification.ID, TimeStamp:notification.TimeStamp, MerchantID:notification.MerchantID, ActivityType:notification.ActivityType, UnRead:notification.UnRead, Message:notification.Message}
		switch notification.ActivityType {
		case "request":
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		}

		notifications = append(notifications, notificationInfo)
	}
	return notifications, err
}

func GetUserNotifications(userID string) (notifications []interface{}, err error) {
	notifications = []interface{}{}
	r, err := DB.Query(`SELECT u.id, ts, user_id, type, source_id, unread, message 
								FROM user_notification u
								LEFT JOIN activity_type a on u.activity_type = a.id
								WHERE user_id=?`, userID)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer r.Close()
	var notification models.UserNotification
	for r.Next()  {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.UserID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
		}

		notificationInfo := models.UserNotificationInfo{ID: notification.ID, TimeStamp:notification.TimeStamp, UserID:notification.UserID, ActivityType:notification.ActivityType, UnRead:notification.UnRead, Message:notification.Message}
		switch notification.ActivityType {
		case "request":
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		}

		notifications = append(notifications, notificationInfo)
	}
	return notifications, err
}

func CreateNotificationByUserId(userID string, title string, message string, claims models.JwtClaims, data interface{}) (res interface{}, err error) {
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

func SendSMSMessage(receiver string, pin string) {
	twilioPhoneNumb := "+17243053011"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/AC83d1ece83480dc998468d333f73e12ed/Messages.json"
	accountSid := "AC83d1ece83480dc998468d333f73e12ed"
	authToken := "e140ee76a6b3b3166ba5490d00faef47"

	msgData := url.Values{}
	msgData.Set("To", receiver)
	msgData.Set("From", twilioPhoneNumb)
	msgData.Set("Body", pin)

	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&data)

		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}
