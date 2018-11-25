package controllers

import (
	"encoding/json"
	"errors"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"net/http"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	if claims.Role != "admin" {
		return errors.New("Failed to authorize")
	}

	switch req.Method {
	case "GET":
		merchant, err := services.GetMerchantByEmail(objectID)
		if err != nil {
			return err
		}
		b, _ := json.Marshal(merchant)
		w.Write([]byte(b))

	default:
		http.NotFound(w, req)
	}

	return nil
}
