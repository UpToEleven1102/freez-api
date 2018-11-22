package controllers

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../services"
	"../models"
	"github.com/dgrijalva/jwt-go"
	"time"
	"os"
)

type _tokenResponse struct {
	token string
}

func createToken(merchant models.Merchant) (string, error) {

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Id:        merchant.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func SignUp(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var merchant models.Merchant

	json.Unmarshal(body, &merchant)

	merchant, err = services.CreateMerchant(merchant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, _ := createToken(merchant)
	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func SignIn(w http.ResponseWriter, req *http.Request) {

}
