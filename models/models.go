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
	Option       merchantMOption `json:"option"`
	Product      []Product       `json:"product"`
}

type merchantMOption struct {
	AddConvenienceFee bool `json:"add_convenience_fee"`
}

type Product struct {
	ID         int     `json:"id"`
	MerchantId string  `json:"merchant_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Image      string  `json:"image"`
}

type MerchantInfo struct {
	Online      bool    `json:"online"`
	MerchantID  string  `json:"merchant_id"`
	Location    LongLat `json:"location"`
	Distance    float32 `json:"distance"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Mobile      bool    `json:"mobile"`
	Image       string  `json:"image"`
	IsFavorite  bool    `json:"is_favorite"`
	Accepted    int     `json:"accepted"`
	//Product string `json:"product"`
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

type RequestInfo struct {
	ID          int     `json:"id"`
	UserId      string  `json:"user_id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Location    LongLat `json:"location"`
	Comment     string  `json:"comment"`
	PhoneNumber string  `json:"phone_number"`
	Image       string  `json:"image"`
	Distance    float32 `json:"distance"`
	Accepted    int     `json:"accepted"`
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
