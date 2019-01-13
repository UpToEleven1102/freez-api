package models

type WSRequestData struct {
	Type string `json:"type"`
	Payload string `json:"payload"`
}

type UserWebSocketRequestData struct {
	Token string `json:"auth_token"`
	Location LongLat `json:"location"`
}
