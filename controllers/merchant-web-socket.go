package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/identity"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const (
	requestInfo  = "request_info"
	postLocation = "post_location"
)

/*MerchantWebSocketHandler - HandleFunc for merchant websocket route

Notes: notify user merchant nearby every 6 hours*/
func MerchantWebSocketHandler(ws *websocket.Conn) {
	var (
		claims       models.JwtClaims
		reqData      models.WSRequestData
		err          error
		claimSt      string
		secSocketKey string
	)

	defer ws.Close()

	for {
		secSocketKey = ws.Request().Header.Get("Sec-WebSocket-Key")

		_ = services.RedisClient.Get(secSocketKey).Scan(&claimSt)
		_ = json.Unmarshal([]byte(claimSt), &claims)

		reqData, err = parseMessage(ws)
		if err != nil {
			log.Println(err)
			break
		}

		switch reqData.Type {
		case token:
			if claims, err = identity.AuthenticateToken(reqData.Payload); err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message: err.Error()})
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}
			}
			b, _ := json.Marshal(claims)
			services.RedisClient.Set(secSocketKey, string(b), time.Hour*24)

		case requestInfo:
			requests, err := services.GetRequestInfoByMerchantId(claims.Id)

			if err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message: err.Error()})
				_ = websocket.Message.Send(ws, string(b))
				break
			}

			b, _ := json.Marshal(requests)
			b, _ = json.Marshal(models.DataResponse{Success: true, Type: requestInfo, Message: string(b)})
			if err = websocket.Message.Send(ws, string(b)); err != nil {
				break
			}
		case postLocation:
			var location models.Location
			if err = json.Unmarshal([]byte(reqData.Payload), &location.Location); err != nil {
				break
			}

			location.Id = claims.Id
			if err = services.AddNewLocation(location); err != nil {
				break
			}

			//push notification to user when the merchant is nearby
			var userIds []interface{}
			if userIds, err = services.GetUserIDNotifyMerchantNearbyByMerchantID(location); err != nil {
				break
			}

			for _, userID := range userIds {
				var merchantID string
				r := services.RedisClient.Get("nearby_notification_" + userID.(string) + "_" + claims.Id)
				if err = r.Scan(&merchantID); err != nil {
					merchantID = claims.Id

					claims.Id = userID.(string)
					claims.Role = config.User
					services.RedisClient.Set("nearby_notification_"+userID.(string)+"_"+merchantID, merchantID, time.Hour*6)
					err = services.CreateNotification(config.NotifTypeFavNearbyID, -1, merchantID, "Favorite Nearby", "Your favorite merchant is nearby", claims)
					if err != nil {
						log.Println(err)
					}
				}
				//do something else
			}
		}

		if err != nil {
			log.Println("merchant socket log err ", err)
			break
		}
	}
}
