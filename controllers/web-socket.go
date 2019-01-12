package controllers

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const (
	merchantNearby = "merchant_nearby"
	token          = "token"
)

func parseMessage(ws *websocket.Conn) (interface{}, error) {
	var message string
	if err := websocket.Message.Receive(ws, &message); err != nil {
		log.Println(err)
		return nil, err
	}

	var reqData models.WSRequestData

	if err := json.Unmarshal([]byte(message), &reqData); err != nil {
		log.Println(err)
		return nil, err
	}

	return reqData, nil
}

var (
	claims  models.JwtClaims
	reqData models.WSRequestData
)

func UserWebSocketHandler(ws *websocket.Conn) {
	var claimSt string
	var secSocketKey string

	for {
		secSocketKey = ws.Request().Header.Get("Sec-WebSocket-Key")
		_ = services.RedisClient.Get(secSocketKey).Scan(&claimSt)
		_ = json.Unmarshal([]byte(claimSt), &claims)

		parsedMsg, err := parseMessage(ws)
		if err != nil {
			break
		}

		reqData = parsedMsg.(models.WSRequestData)

		switch reqData.Type {
		case token:
			claims, err = identity.AuthenticateToken(reqData.Payload)
			if err != nil {
				if err := websocket.Message.Send(ws, models.DataResponse{Success: false}); err != nil {
					log.Println(err)
					break
				}
			}
			b, _ := json.Marshal(claims)
			services.RedisClient.Set(secSocketKey, string(b), time.Hour*24)

		case merchantNearby:
			_, _ = getMerchantNearby()
		}

		response := "received message"

		if err = websocket.Message.Send(ws, response); err != nil {
			log.Println(err)
			break
		}
	}
}

func getMerchantNearby() (merchants []interface{}, err error) {
	var location models.Location
	location.Id = claims.Id

	_ = json.Unmarshal([]byte(reqData.Payload), &location.Location)

	fmt.Printf("%+v", location)

	return nil, nil
}
