package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"github.com/stripe/stripe-go"
	"log"
	"net/http"
)

type RequestObject struct {
	StripeID string `json:"stripe_id"`
	CardID   string `json:"card_id"`
	Token    string `json:"token"`
}

func StripeOpsHandler(w http.ResponseWriter, req *http.Request, urlString string, claims models.JwtClaims) error {
	if claims.Role != config.Merchant {
		return errors.New("Unauthorized")
	}

	objectID, _ := getUrlParam(urlString)

	var jsonEncoder = json.NewEncoder(w)

	switch req.Method {
	case "GET":
		switch objectID {
		case "card-list":
			cards, err := services.GetStripeCardList(claims.Id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
			}

			_ = json.NewEncoder(w).Encode(cards)

		case "account":
			acc, err := services.GetMerchantStripeAccount(claims.Id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(acc)

		case "acc-balance":
			id, err := services.GetMerchantStripeIdByMerchantId(claims.Id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			b, err := services.StripeGetAccountBalance(id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(b)

		default:
			http.NotFound(w, req)
		}

	case "POST":
		switch objectID {
		case "card":
			var data RequestObject

			_ = json.NewDecoder(req.Body).Decode(&data)
			fmt.Println(data)

			c, err := services.StripeConnectCreateDebitCard(data.StripeID, data.Token)

			if err != nil {
				log.Println(err.Error())
				_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = jsonEncoder.Encode(c)

		case "refund":
			type refundData struct {
				OrderID  int     `json:"order_id"`
				StripeID string  `json:"stripe_id"`
				Amount   float64 `json:"amount"`
				Reason   string  `json:"reason"`
			}

			var data refundData

			err := json.NewDecoder(req.Body).Decode(&data)

			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			var res interface{}

			if data.Amount <= 0 {
				res, err = services.StripeRefund(data.StripeID)
			} else {
				res, err = services.StripePartialRefund(data.StripeID, data.Amount)
			}

			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			refundRes := res.(*stripe.Refund)

			fmt.Printf("%+v \n", res)
			//_ = json.NewEncoder(w).Encode(res)

			orderEntity, err := services.GetOrderEntityById(data.OrderID)
			order := orderEntity.(models.OrderEntity)
			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			var userClaims models.JwtClaims
			userClaims.Id = order.UserId
			userClaims.Role = config.User
			err = services.CreateNotification(config.NotifTypeRefundMadeID, int64(order.ID), order.MerchantId, "Order refunded", "Refund for your order started", userClaims )

			if err != nil {
				panic(json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()}))
			}

			_ = json.NewEncoder(w).Encode(models.DataResponse{Success: refundRes.Status == "succeeded", Message: "Successfully refunded"})

		default:
			http.NotFound(w, req)
		}

	case "DELETE":
		switch objectID {
		case "card":
			var data RequestObject

			_ = json.NewDecoder(req.Body).Decode(&data)
			c, err := services.StripeConnectDeleteDebitCard(data.StripeID, data.CardID)

			if err != nil {
				log.Println(err)
				_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = jsonEncoder.Encode(c)
		default:
			http.NotFound(w, req)
		}

	case "PUT":
		switch objectID {
		case "card":
			var data RequestObject

			_ = json.NewDecoder(req.Body).Decode(&data)
			fmt.Println(data)

			card, err := services.StripeConnectMakeDefaultCurrencyDebitCard(data.StripeID, data.CardID)

			if err != nil {
				log.Println(err)
				_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = jsonEncoder.Encode(card)

		case "verify-entity":
			var data models.MerchantEntityInformation
			_ = json.NewDecoder(req.Body).Decode(&data)

			fmt.Println(data)

			_, err := services.StripeConnectEntityVerification(claims.Id, data)

			if err != nil {
				log.Println(err)

				_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = jsonEncoder.Encode(models.DataResponse{Success: true})

		default:
			http.NotFound(w, req)
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
