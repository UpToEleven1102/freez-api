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

var (
	claims  models.JwtClaims
	reqData models.WSRequestData
	err error
)

func parseMessage(ws *websocket.Conn) (models.WSRequestData, error) {
	var message string
	if err := websocket.Message.Receive(ws, &message); err != nil {
		return reqData, err
	}

	if err := json.Unmarshal([]byte(message), &reqData); err != nil {
		return reqData, err
	}

	return reqData, nil
}



func UserWebSocketHandler(ws *websocket.Conn) {
	var claimSt string
	var secSocketKey string

	defer ws.Close()

	for {
		secSocketKey = ws.Request().Header.Get("Sec-WebSocket-Key")
		_ = services.RedisClient.Get(secSocketKey).Scan(&claimSt)
		_ = json.Unmarshal([]byte(claimSt), &claims)

		reqData, err := parseMessage(ws)
		if err != nil {
			break
		}

		switch reqData.Type {
		case token:
			claims, err = identity.AuthenticateToken(reqData.Payload)
			if err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message:err.Error()})
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}
			}
			b, _ := json.Marshal(claims)
			services.RedisClient.Set(secSocketKey, string(b), time.Hour*24)

		case merchantNearby:
			var merchants []interface{}
			merchants, err = getMerchantNearby()

			if err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message:err.Error()})
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}
			}
			b, _ := json.Marshal(merchants)
			b, _ = json.Marshal(models.DataResponse{Success:true, Type: reqData.Type, Message: string(b)})

			err = websocket.Message.Send(ws,string(b))
			if err != nil {
				break
			}

		case postLocation:
			var user models.User
			user.ID = claims.Id

			if err = json.Unmarshal([]byte(reqData.Payload), &user.LastLocation); err != nil {
				break
			}

			fmt.Println(user)

			if _ , err = services.UpdateUserLocation(user); err != nil {
				break
			}
		}

		if err != nil {
			log.Println(err)
			break
		}
	}
}

func getMerchantNearby() (merchants []interface{}, err error) {
	var location models.Location
	location.Id = claims.Id

	_ = json.Unmarshal([]byte(reqData.Payload), &location.Location)
	return services.GetNearMerchantsLastLocation(location)
}
