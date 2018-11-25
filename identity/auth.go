package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	c "git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type JWTData struct {
	jwt.StandardClaims
	//CustomClaims map[string]string `json:"custom,omitempty"`
}

type response struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func getToken(w http.ResponseWriter, req *http.Request) (string, error) {
	authToken := req.Header.Get("Authorization")
	authArr := strings.SplitN(authToken, " ", 2)

	if len(authArr) != 2 {
		return "", errors.New("Authentication header is invalid" + authToken)
	}
	return authArr[1], nil
}

func AuthorizeMiddleware(w http.ResponseWriter, req *http.Request, objectID string, handler models.FuncHandler) error {
	token, err := getToken(w, req)
	if err != nil{
		return err
	}
	claims, err := AuthenticateToken(token)
	if err != nil {
		return err
	}

	err = handler(w, req, objectID, claims)
	return err
}


func AuthenticateToken(jwtToken string) (models.JwtClaims, error) {
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

		return models.JwtClaims{Id: id, Role: role}, err
	}
	return models.JwtClaims{}, err
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
	err := json.Unmarshal(body, &r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	token, err := getToken(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	claims, err := AuthenticateToken(token)

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
