package models

type Merchant struct {
	ID           string          `json:"id"`
	Online       bool            `json:"online"`
	PhoneNumber  string          `json:"phone_number"`
	Email        string          `json:"email"`
	Name         string          `json:"name"`
	Mobile       bool            `json:"mobile"`
	Password     string          `json:"password"`
	Image        string          `json:"image"`
	Role         string          `json:"role"`
	LastLocation LongLat         `json:"last_location"`
	StripeID     string          `json:"stripe_id"`
	CardID       string          `json:"card_id"`
	FacebookID   string          `json:"facebook_id"`
	Category     string          `json:"category"`
	Option       merchantMOption `json:"option"`
	Product      []Product       `json:"product"`
}

type merchantMOption struct {
	AddConvenienceFee bool `json:"add_convenience_fee"`
}

type MerchantNotification struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	TimeStamp    string `json:"ts"`
	MerchantID   string `json:"merchant_id"`
	ActivityType string `json:"activity_type"`
	SourceID     int    `json:"source_id"`
	UnRead       bool   `json:"unread"`
	Message      string `json:"message"`
}

type MerchantNotificationInfo struct {
	ID           int         `json:"id"`
	TimeStamp    string      `json:"ts"`
	MerchantID   string      `json:"merchant_id"`
	ActivityType string      `json:"activity_type"`
	Source       interface{} `json:"source"`
	UnRead       bool        `json:"unread"`
	Message      string      `json:"message"`
}

type MerchantInfo struct {
	Online      bool          `json:"online"`
	Category    string        `json:"category"`
	MerchantID  string        `json:"merchant_id"`
	Location    LongLat       `json:"location"`
	Distance    float32       `json:"distance"`
	Name        string        `json:"name"`
	PhoneNumber string        `json:"phone_number"`
	Email       string        `json:"email"`
	Mobile      bool          `json:"mobile"`
	Image       string        `json:"image"`
	IsFavorite  bool          `json:"is_favorite"`
	Accepted    int           `json:"accepted"`
	Products    []interface{} `json:"products"`
}

type Product struct {
	ID         int     `json:"id"`
	MerchantId string  `json:"merchant_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Image      string  `json:"image"`
}
