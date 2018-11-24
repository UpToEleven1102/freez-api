package identity

import (
	"encoding/json"
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

	json.Unmarshal(body, &credentials)

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
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var merchant models.Merchant

	json.Unmarshal(body, &merchant)

	merchant, err = services.CreateMerchant(merchant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, _ := createToken(merchant)
	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
