package services

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/tbalthazar/onesignal-go"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

func CreateNotification(activityType int, sourceID int64, id string, title string, message string, claims models.JwtClaims) (err error) {
	switch claims.Role {
	case config.User:
		err := InsertUserNotification(claims.Id, activityType, sourceID, id, message)

		if err != nil {
			panic(err)
			return err
		}

	case config.Merchant:
		err := InsertMerchantNotification(claims.Id, activityType, sourceID, message)

		if err != nil {
			panic(err)
			return err
		}
	}

	var data interface{}
	_, err = SendNotificationByUserId(claims.Id, title, message, data)

	if err != nil {
		panic(err)
		return err
	}

	return nil

}

func SendNotificationByUserId(userID string, title string, message string, data interface{}) (res interface{}, err error) {
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
