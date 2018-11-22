package models


type Merchant struct {
	ID string `json:"id"`
	PhoneNumber string `json:"phone-number"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Image string `json:"image"`
}
