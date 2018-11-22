package controllers

import (
	"net/http"
	"encoding/json"
	"../services"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	method := req.Method
	switch method {
	case "GET":
		merchant, err := services.GetMerchantByEmail(objectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		b, _ := json.Marshal(merchant)
		w.Write([]byte(b))

	default:
		http.NotFound(w, req)
	}
}
