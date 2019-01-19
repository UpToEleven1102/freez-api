package config

const (
	SignUp = "signup"
	SignIn = "signin"
	Merchant = "merchant"
	User = "user"
	UserInfo = "userinfo"
	Exist = "exist"
	Verify = "verify"
	Request = "request"
	Location = "location"
	Stripe = "stripe"
	Email = "email"
	PhoneNumber="phone-number"
	SRID = 3857

	NOTIF_TYPE_FLAG_REQUEST = "request"
	NOTIF_TYPE_PAYMENT_MADE = "payment_made"
	NOTIF_TYPE_REFUND_MADE  = "refund_made"
	NOTIF_TYPE_FAV_NEARBY   = "merchant_nearby"

	NOTIF_TYPE_FLAG_REQUEST_ID = 1
	NOTIF_TYPE_PAYMENT_MADE_ID = 2
	NOTIF_TYPE_REFUND_MADE_ID  = 3
	NOTIF_TYPE_FAV_NEARBY_ID   = 4
)
