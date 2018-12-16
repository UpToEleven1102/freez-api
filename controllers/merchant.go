package controllers

import (
	"encoding/json"
	"errors"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"net/http"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	switch req.Method {
	case "GET":
		switch objectID {
		case "":
			if claims.Role != "admin" {
				return errors.New("Failed to authorize")
			}
			merchant, err := services.GetMerchantByEmail(objectID)
			if err != nil {
				return err
			}
			b, _ := json.Marshal(merchant)
			w.Write([]byte(b))
		}

	case "POST":
		if objectID == "update-status" {
			id := claims.Id
			err := services.ChangeOnlineStatus(id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
		} else {
			http.NotFound(w, req)
			return nil
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
