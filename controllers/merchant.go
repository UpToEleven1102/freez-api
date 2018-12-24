package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
	"strings"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	var response models.DataResponse

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
			_,_ = w.Write(b)
		case "presign-url":
			fileName := fmt.Sprint(claims.Id, "-profile.jpg")
			url, err := services.GeneratePreSignedUrl(fileName)

			var response models.DataResponse

			if err != nil {
				response.Success = false
				response.Message = err.Error()
			} else {
				response.Success = true
				response.Message = url
			}

			b, _ := json.Marshal(response)
			_,_ = w.Write(b)
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

	case "PUT":
		if objectID == "update-profile" {
			b, err := ioutil.ReadAll(req.Body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			var merchant models.Merchant
			err = json.Unmarshal(b, &merchant)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
			merchant.ID = claims.Id

			err = services.UpdateMerchant(merchant)

			if err != nil {
				response.Success = false
				if strings.Contains(err.Error(), "Error 1062") {
					response.Message = "Email is currently in use!"
				} else {
					response.Message = err.Error()
				}

				sendResponse(w, response, http.StatusBadRequest)
			}
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
