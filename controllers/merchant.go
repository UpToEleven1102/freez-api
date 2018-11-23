package controllers

import (
	"net/http"
	"encoding/json"
	"../services"
	auth "../identity"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	if claims, err := auth.AuthenticateTokenMiddleWare(w,req); err!=nil && claims.Role!="admin" {
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
