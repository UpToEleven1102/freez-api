package identity

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

func SignInMerchant(w http.ResponseWriter, req *http.Request) {
	var credentials Credentials
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := services.GetMerchantByEmail(credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if res != nil {
		merchant := res.(models.Merchant)
		err = bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(credentials.Password))
		if err == nil {
			token, err := createToken(merchant)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			b, _ := json.Marshal(token)
			w.WriteHeader(http.StatusAccepted)
			w.Write(b)
			return
		}
	}
	http.Error(w, "Credentials Invalid", http.StatusBadRequest)
}

func SignUpMerchant(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	var response models.DataResponse
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusBadRequest)
		fmt.Println(response)
		return
	}

	var merchant models.Merchant

	err = json.Unmarshal(body, &merchant)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusBadRequest)
		return
	}

	merchant, err = services.CreateMerchant(merchant)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusInternalServerError)
		fmt.Println(response)
		return
	}

	token, _ := createToken(merchant)
	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
