package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freez-app-rest/identity"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const (
	merchantNearby = "merchant_nearby"
	token          = "token"
)

type requestMerchantData struct {
	Location models.LongLat `json:"location"`
	Filter   []string       `json:"filter"`
}

func parseMessage(ws *websocket.Conn) (reqData models.WSRequestData, err error) {
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
			claims, err = identity.AuthenticateToken(reqData.Payload)
			if err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message: err.Error()})
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}
			}
			b, _ := json.Marshal(claims)
			services.RedisClient.Set(secSocketKey, string(b), time.Hour*24)

		case merchantNearby:
			var merchants []interface{}
			merchants, err = getMerchantNearby(claims, reqData)

			if err != nil {
				b, _ := json.Marshal(models.DataResponse{Success: false, Message: err.Error()})
				if err = websocket.Message.Send(ws, string(b)); err != nil {
					break
				}
			}
			b, _ := json.Marshal(merchants)
			b, _ = json.Marshal(models.DataResponse{Success: true, Type: reqData.Type, Message: string(b)})

			err = websocket.Message.Send(ws, string(b))
			if err != nil {
				break
			}

		case searchMerchant:
			var searchData models.SearchData
			var location models.Location

			err := json.Unmarshal([]byte(reqData.Payload), searchData)

			if err != nil {
				log.Println(err.Error())
				break
			}

			err = json.Unmarshal([]byte(reqData.Payload), location)

			if err != nil {
				log.Println(err.Error())
				break
			}


			log.Printf("%+v", searchData)
			log.Printf("%+v", location)

			merchants, err := services.FilterMerchantByName(searchData, location)

			var b []byte

			if err != nil {
				b, _ = json.Marshal(models.DataResponse{Success:false, Message:err.Error()})
			} else {
				b, _ = json.Marshal(merchants)
				b, _ = json.Marshal(models.DataResponse{Success:true, Type: merchantNearby, Message: string(b)})
			}

			if err = websocket.Message.Send(ws, string(b)); err != nil {
				break
			}

		case postLocation:
			var user models.User
			user.ID = claims.Id

			if err = json.Unmarshal([]byte(reqData.Payload), &user.LastLocation); err != nil {
				break
			}

			if _, err = services.UpdateUserLocation(user); err != nil {
				break
			}
		}

		if err != nil {
			log.Println(err)
			break
		}
	}
}

func getMerchantNearby(claims models.JwtClaims, reqData models.WSRequestData) (merchants []interface{}, err error) {
	var data requestMerchantData
	var location models.Location

	_ = json.Unmarshal([]byte(reqData.Payload), &data)

	location.Location = data.Location
	location.Id = claims.Id
	return services.GetNearbyMerchantsLastLocation(location, data.Filter...)
}
