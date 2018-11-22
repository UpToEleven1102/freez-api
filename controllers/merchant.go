package controllers

import (
	"net/http"
	"encoding/json"
	"github.com/UpToEleven1102/freezeapp-rest/services"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, attribute string, objectID string) {
	method := req.Method
	switch method {
	case "GET":
		if len(attribute) == 0 {
			if len(objectID) != 0 {
				http.NotFound(w, req)
			}
			merchants, err := services.GetMerchants()

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			b, _ := json.Marshal(merchants)
			w.Write([]byte(b))
		} else {
			if len(objectID) == 0 {
				http.NotFound(w, req)
			}

			switch attribute {
			case "email":
				merchant, err := services.GetMerchantByEmail(objectID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				b, _ := json.Marshal(merchant)
				w.Write([]byte(b))
			}
		}
	default:
		http.NotFound(w, req)
	}
}
