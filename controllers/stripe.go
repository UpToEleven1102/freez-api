package controllers

import (
	"encoding/json"
	"errors"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"log"
	"net/http"
)

type RequestObject struct {
	StripeID string `json:"stripe_id"`
	CardID string `json:"card_id"`
}

func StripeOpsHandler(w http.ResponseWriter, req *http.Request, urlString string, claims models.JwtClaims) error {
	if claims.Role != config.Merchant {
		return errors.New("Unauthorized")
	}

	objectID, _ := getUrlParam(urlString)

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
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(acc)

		case "acc-balance":
			id, err := services.GetMerchantStripeIdByMerchantId(claims.Id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()})
				return nil
			}

			b, err := services.StripeGetAccountBalance(id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(b)

		default:
			http.NotFound(w, req)
		}
	case "DELETE":
		switch objectID {
		case "card":
			var data RequestObject

			var jsonEncoder = json.NewEncoder(w)
			_ = json.NewDecoder(req.Body).Decode(&data)
			c, err := services.StripeConnectDeleteDebitCard(data.StripeID, data.CardID)

			if err != nil {
				log.Println(err)
				_ = jsonEncoder.Encode(models.DataResponse{Success:false, Message:err.Error()})
				return nil
			}

			_ = jsonEncoder.Encode(c)
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
