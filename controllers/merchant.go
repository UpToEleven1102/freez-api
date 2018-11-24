package controllers

import (
	"encoding/json"
	auth "git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"net/http"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	if claims, err := auth.AuthenticateTokenMiddleWare(w, req); err != nil && claims.Role != "admin" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	method := req.Method
	switch method {
	case "GET":
		merchant, err := services.GetMerchantByEmail(objectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b, _ := json.Marshal(merchant)
		w.Write([]byte(b))

	default:
		http.NotFound(w, req)
	}
}
