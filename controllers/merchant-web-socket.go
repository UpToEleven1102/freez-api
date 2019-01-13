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
	requestInfo = "request_info"
)

func MerchantWebSocketHandler(ws *websocket.Conn) {
	var claimSt string
	var secSocketKey string

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
				b, _ := json.Marshal(models.DataResponse{Success: false, Message:err.Error()})
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
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}

				break
			}
			fmt.Println(claims)

			b, _ := json.Marshal(requests)
			b, _ = json.Marshal(models.DataResponse{Success: true, Type: requestInfo, Message: string(b)})
			if err = websocket.Message.Send(ws, string(b) ); err != nil {
				break
			}
		}

		if err != nil {
			log.Println(err)
			break
		}
	}
}
