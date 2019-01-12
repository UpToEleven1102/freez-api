package models

import (
	"net/http"
)

type AuthFuncHandler func(http.ResponseWriter, *http.Request, string, string)

type FuncHandler func(http.ResponseWriter, *http.Request, string, JwtClaims) error

type JwtClaims struct {
	Id   string
	Role string
}

type RequestData struct {
	UserId string `json:"user_id"`
	Data   string `json:"data"`
}

type ItemHistory struct {
	ID       int     `json:"id"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderRequestData struct {
	StripeToken string  `json:"stripe_token"`
	Amount      float64 `json:"amount"`
	UserID      string  `json:"user_id"`
	MerchantID  string  `json:"merchant_id"`
	StripeID    string  `json:"stripe_id"`
	Items       []struct {
		Product  Product `json:"product"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"items"`
}

type OrderEntity struct {
	ID         int         `json:"id"`
	UserId     string      `json:"user_id"`
	MerchantId string      `json:"merchant_id"`
	StripeId   string      `json:"stripe_id"`
	Refund     bool        `json:"refund"`
	Amount     float64     `json:"amount"`
	Date       string      `json:"date"`
	Items      interface{} `json:"items"`
}

type RequestInfo struct {
	ID          int     `json:"id"`
	UserId      string  `json:"user_id"`
	MerchantId  string  `json:"merchant_id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Location    LongLat `json:"location"`
	Comment     string  `json:"comment"`
	PhoneNumber string  `json:"phone_number"`
	Image       string  `json:"image"`
	Distance    float32 `json:"distance"`
	Accepted    int     `json:"accepted"`
	Active      bool    `json:"active"`
}

type RequestNotification struct {
	ID       int          `json:"id"`
	User     User         `json:"user"`
	Merchant MerchantInfo `json:"merchant"`
	Location LongLat      `json:"location"`
	Comment  string       `json:"comment"`
	Active   bool         `json:"active"`
	Accepted int          `json:"accepted"`
}

type Request struct {
	UserId     string  `json:"user_id"`
	Comment    string  `json:"comment"`
	MerchantID string  `json:"merchant_id"`
	Location   LongLat `json:"location"`
}

type RequestEntity struct {
	ID         int     `json:"id"`
	UserID     string  `json:"user_id"`
	MerchantID string  `json:"merchant_id"`
	Location   LongLat `json:"location"`
	Comment    string  `json:"comment"`
	Active     bool    `json:"active"`
	Accepted   int     `json:"accepted"`
}

type LongLat struct {
	Long float32 `json:"long"`
	Lat  float32 `json:"lat"`
}

type Location struct {
	Id       string  `json:"id"`
	Location LongLat `json:"location"`
}

type LocationEntity struct {
	ID         int     `json:"id"`
	MerchantID string  `json:"merchant_id"`
	Location   LongLat `json:"location"`
}

type DataResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
