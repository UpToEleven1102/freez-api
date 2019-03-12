package services

import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"log"
)


func InsertUserNotification(userID string, activityType int, sourceID int64, merchantID string, message string) error {
	_, err := DB.Exec(`INSERT INTO user_notification (user_id, activity_type, source_id, merchant_id, message) VALUES (?, ?, ?, ?, ?)`, userID, activityType, sourceID, merchantID, message)
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

func GetUserNotifications(userID string) (notifications []interface{}, err error) {
	notifications = []interface{}{}
	r, err := DB.Query(`SELECT u.id, ts, user_id, type, source_id, merchant_id, unread, message 
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
		err = r.Scan(&notification.ID, &notification.TimeStamp, &notification.UserID, &notification.ActivityType, &notification.SourceID, &notification.MerchantID, &notification.UnRead, &notification.Message)
		if err != nil {
			log.Println(err)
		}

		notificationInfo := models.UserNotificationInfo{ID: notification.ID, TimeStamp:notification.TimeStamp, UserID:notification.UserID, ActivityType:notification.ActivityType, UnRead:notification.UnRead, Message:notification.Message}
		switch notification.ActivityType {
		case config.NOTIF_TYPE_FLAG_REQUEST:
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		case config.NOTIF_TYPE_FAV_NEARBY:
			notificationInfo.Merchant = models.MerchantInfo{MerchantID:notification.MerchantID}
		case config.NOTIF_TYPE_REFUND_MADE:
			notificationInfo.Source, _ = GetOrderById(notification.SourceID)
		}

		log.Println(notificationInfo)

		notifications = append(notifications, notificationInfo)
	}
	return notifications, err
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
		case config.NOTIF_TYPE_FLAG_REQUEST:
			notificationInfo.Source, _ = GetRequestNotificationById(notification.SourceID)
		case config.NOTIF_TYPE_REFUND_MADE:
			notificationInfo.Source, _ = GetOrderById(notification.SourceID)
		case config.NOTIF_TYPE_FAV_NEARBY:
			notificationInfo.Merchant = models.MerchantInfo{MerchantID:notification.MerchantID}
		}
		return notificationInfo, nil
	}
	return nil, nil
}
