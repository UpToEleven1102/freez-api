package identity

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../models"
	"../services"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var user models.User
	json.Unmarshal(body, &user)

	r, err := services.CreateUser(user)
	if err != nil {
		http.Error(w, "account exists", http.StatusBadRequest)
		return
	}

	user = r.(models.User)
	token, err := createToken(user)

	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func SignInUser(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var credentials Credentials
	json.Unmarshal(body, &credentials)

	r, err := services.GetUserByEmail(credentials.Email)

	if err != nil || r == nil {
		http.Error(w, "Credentials invalid", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(r.(models.User).Password), []byte(credentials.Password))

	if err != nil {
		http.Error(w, "Credentials invalid", http.StatusBadRequest)
		return
	}

	token, _ := createToken(r)
	b, _:= json.Marshal(token)
	w.Write(b)
}
