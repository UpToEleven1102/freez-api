package models

type UserNotification struct {
	ID           int    `json:"id"`
	TimeStamp    string `json:"ts"`
	UserID       string `json:"user_id"`
	ActivityType string `json:"activity_type"`
	SourceID     int    `json:"source_id"`
	UnRead       bool   `json:"unread"`
	Message      string `json:"message"`
}

type UserNotificationInfo struct {
	ID           int         `json:"id"`
	TimeStamp    string      `json:"ts"`
	UserID       string      `json:"user_id"`
	ActivityType string      `json:"activity_type"`
	Source       interface{} `json:"source"`
	UnRead       bool        `json:"unread"`
	Message      string      `json:"message"`
}

type User struct {
	ID           string  `json:"id"`
	PhoneNumber  string  `json:"phone_number"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Password     string  `json:"password"`
	LastLocation LongLat `json:"last_location"`
	Image        string  `json:"image"`
	Option       mOption `json:"option"`
	Role         string  `json:"role"`
}

type mOption struct {
	NotifFavNearby bool `json:"notif_fav_nearby"`
}