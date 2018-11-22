package controllers

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../services"
	"../models"
	c "../config"
	"github.com/dgrijalva/jwt-go"
	"os"
	"fmt"
	"strings"
	"log"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type response struct {
	Message string `json:"message"`
}

type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type JWTData struct {
	jwt.StandardClaims
	//CustomClaims map[string]string `json:"custom,omitempty"`
}

type JwtClaims struct {
	Id   string
	Role string
}

func AuthenticateTokenMiddleWare(w http.ResponseWriter, req *http.Request) (JwtClaims, error) {
	authToken := req.Header.Get("Authorization")
	authArr := strings.SplitN(authToken, " ", 2)

	if len(authArr) != 2 {
		log.Println("Authentication header is invalid" + authToken)
		return JwtClaims{}, errors.New("Request failed!")
	}
	jwtToken := authArr[1]
	token, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if claims, ok := token.Claims.(*JWTData); ok && token.Valid {
		err = claims.Valid()

		id := claims.Id
		role := claims.Subject

		return JwtClaims{Id: id, Role: role}, err
	}
	return JwtClaims{}, err
}

func createToken(merchant models.Merchant) (string, error) {
	claims := JWTData{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Id:        merchant.ID,
			IssuedAt:  time.Now().Unix(),
			Subject:   "merchant",
		},
		//map[string]string{
		//	"Id": merchant.ID,
		//	"Role": "merchant",
		//},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	secretKey := os.Getenv("SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func merchantExists(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var r response
	json.Unmarshal(body, &r)

	if	merchant, _ := services.GetMerchantByEmail(r.Message); merchant!= nil {
		r.Message = "true"
	} else {
		r.Message = "false"
	}

	b, _ := json.Marshal(r)
	w.Write(b)
}

func signUpMerchant(w http.ResponseWriter, req *http.Request) {
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

func getUserInfo(w http.ResponseWriter, req *http.Request) {
	claims, err := AuthenticateTokenMiddleWare(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	if claims.Role == c.Merchant {
		merchant, err := services.GetMerchantById(claims.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if merchant == nil {
			http.Error(w, "User not exists", http.StatusUnauthorized)
			return
		}
		res := merchant.(models.Merchant)
		res.Role = claims.Role
		b, _ := json.Marshal(res)
		w.Write(b)

	} else if claims.Role == c.User {
		//do something
	}
}

func signInMerchant(w http.ResponseWriter, req *http.Request) {
	var credentials Credentials
	body, err:= ioutil.ReadAll(req.Body)
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
		err = bcrypt.CompareHashAndPassword([]byte(merchant.Password),[]byte(credentials.Password))
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

func signUp(w http.ResponseWriter, req *http.Request) {

}

func signIn(w http.ResponseWriter, req *http.Request) {

}

func AuthHandler(w http.ResponseWriter, req *http.Request, route string, userType string) {
	switch route {
	case c.SignUp:
		if userType == c.Merchant {
			signUpMerchant(w, req)
		} else if userType == c.User {
			signUp(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.SignIn:
		if userType == c.Merchant {
			signInMerchant(w, req)
		} else if userType == c.User {
			signIn(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.UserInfo:
		if len(userType) > 0 {
			http.NotFound(w, req)
		}
		getUserInfo(w, req)
	case c.Verify:
		if userType == c.Merchant {
			merchantExists(w, req)
		} else if userType == c.User {

		} else {

		}
	default:
		http.NotFound(w, req)
	}
}
