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

type Merchant struct {
	ID          string `json:"id"`
	Online      bool   `json:"online"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Mobile      bool   `json:"mobile"`
	Password    string `json:"password"`
	Image       string `json:"image"`
	Role        string `json:"role"`
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
	//Product string `json:"product"`
}

type User struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Image       string `json:"image"`
	Role        string `json:"role"`
}

type Request struct {
	UserId     string  `json:"user_id"`
	MerchantID string  `json:"merchant_id"`
	Location   LongLat `json:"location"`
}

type RequestEntity struct {
	ID     int    `json:"id"`
	UserId string `json:"user_id"`
	//Location *geos.Geometry `json:"location"`
}

type LongLat struct {
	Long float32 `json:"long"`
	Lat  float32 `json:"lat"`
}

type Location struct {
	MerchantID string  `json:"merchant_id"`
	Location   LongLat `json:"location"`
}

type LocationEntity struct {
	ID         int    `json:"id"`
	MerchantID string `json:"merchant_id"`
	//Location *geos.Geometry `json:"location"`
}
