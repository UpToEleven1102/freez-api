package controllers

import (
	"net/http"
	"encoding/json"
	"../services"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request ,objectID string) {
	method := req.Method
	switch method {
	case "GET":
		if len(objectID) == 0 {
			merchants, err := services.GetMerchants()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			b, _ := json.Marshal(merchants)
			w.Write([]byte(b))
		}

	case "POST":
	}

}
