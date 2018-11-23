package identity

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"log"
	"errors"
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"encoding/json"
	"../services"
	"../models"
	c "../config"
)

type response struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

type Credentials struct {
	Email    string `json:"email"`
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

func createToken(acc interface{}) (string, error) {
	var claims JWTData
	switch acc.(type) {
	case models.Merchant:
		merchant := acc.(models.Merchant)
		claims = JWTData{
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
	case models.User:
		user := acc.(models.User)
		claims = JWTData{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				Id:        user.ID,
				IssuedAt:  time.Now().Unix(),
				Subject:   "user",
			},
			//map[string]string{
			//	"Id": user.ID,
			//	"Role": "user",
			//},
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	secretKey := os.Getenv("SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		panic(err)
	}

	return tokenString, nil
}

func AccountExists(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var r response
	json.Unmarshal(body, &r)

	if merchant, _ := services.GetMerchantByEmail(r.Message); merchant != nil {
		r.Message = "true"
		r.Role = c.Merchant
	} else {
		if merchant, _ = services.GetUserByEmail(r.Message); merchant != nil {
			r.Message = "true"
			r.Role = c.User
		} else {
			r.Message = "false"
		}
	}

	b, _ := json.Marshal(r)
	w.Write(b)
}

func GetUserInfo(w http.ResponseWriter, req *http.Request) {
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
			http.Error(w, "Account not exists", http.StatusUnauthorized)
			return
		}
		res := merchant.(models.Merchant)
		res.Role = claims.Role
		b, _ := json.Marshal(res)
		w.Write(b)

	} else if claims.Role == c.User {
		user, err := services.GetUserById(claims.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user == nil {
			http.Error(w, "Account not exists", http.StatusUnauthorized)
			return
		}
		res := user.(models.User)
		res.Role = claims.Role
		b, _ := json.Marshal(res)
		w.Write(b)
	}
}
