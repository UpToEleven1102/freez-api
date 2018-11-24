package models


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
	ID int `json:"id"`
	Email string `json:"email"`
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
}
