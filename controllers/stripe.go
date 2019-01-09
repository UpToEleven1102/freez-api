package controllers

import (
	"encoding/json"
	"errors"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"net/http"
)

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

		default:
			http.NotFound(w, req)
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
