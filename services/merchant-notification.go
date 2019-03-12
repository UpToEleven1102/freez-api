package services

import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"log"
)

func InsertMerchantNotification(merchantID string, activityType int, sourceID int64, message string) error {
	_, err := DB.Exec(`INSERT INTO merchant_notification (merchant_id, activity_type, source_id, message) VALUES (?, ?, ?, ?)`, merchantID, activityType, sourceID, message)
	if err != nil {
		log.Println(err)
	}
	return err
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
	if r.Next() {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.MerchantID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		notificationInfo := models.MerchantNotificationInfo{ID: notification.ID, TimeStamp: notification.TimeStamp, MerchantID: notification.MerchantID, ActivityType: notification.ActivityType, UnRead: notification.UnRead, Message: notification.Message}
		switch notification.ActivityType {
		case config.NOTIF_TYPE_FLAG_REQUEST:
			notificationInfo.Source, _ = GetRequestInfoById(notification.SourceID, notification.MerchantID)

		case config.NOTIF_TYPE_PAYMENT_MADE:
			notificationInfo.Source, _ = GetOrderPaymentById(notification.SourceID)
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
	for r.Next() {
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.MerchantID, &notification.ActivityType, &notification.SourceID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
		}

		notificationInfo := models.MerchantNotificationInfo{ID: notification.ID, TimeStamp: notification.TimeStamp, MerchantID: notification.MerchantID, ActivityType: notification.ActivityType, UnRead: notification.UnRead, Message: notification.Message}
		switch notification.ActivityType {
		case config.NOTIF_TYPE_FLAG_REQUEST:
			notificationInfo.Source, _ = GetRequestInfoById(notification.SourceID, notification.MerchantID)

		case config.NOTIF_TYPE_PAYMENT_MADE:
			notificationInfo.Source, _ = GetOrderPaymentById(notification.SourceID)
		}

		notifications = append(notifications, notificationInfo)
	}
	return notifications, err
}

func UpdateMerchantNotification(notification models.MerchantNotification) error {
	_, err := DB.Exec(`UPDATE merchant_notification SET unread=? WHERE id=?`, notification.UnRead, notification.ID)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
