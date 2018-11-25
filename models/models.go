package models

import (
	"github.com/paulsmith/gogeos/geos"
	"net/http"
)

type AuthFuncHandler func(http.ResponseWriter,*http.Request, string, string)

type FuncHandler func(http.ResponseWriter, *http.Request, string, JwtClaims) error

type JwtClaims struct {
	Id   string
	Role string
}

type Merchant struct {
	ID string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Image string `json:"image"`
	Role string `json:"role"`
}

type User struct {
	ID string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Image string `json:"image"`
	Role string `json:"role"`
}

type Request struct {
	Email string `json:"email"`
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
}

type RequestEntity struct {
	ID int `json:"id"`
	UserId string `json:"user_id"`
	Location *geos.Geometry `json:"location"`
}

type LatLong struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
}

type Location struct {
	MerchantID string `json:"merchant_id"`
	Location LatLong `json:"location"`
}

type LocationEntity struct {
	ID int `json:"id"`
	MerchantID string `json:"merchant_id"`
	Location *geos.Geometry `json:"location"`
}
