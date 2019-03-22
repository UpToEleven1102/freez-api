package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	c "git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type JWTData struct {
	jwt.StandardClaims
	//CustomClaims map[string]string `json:"custom,omitempty"`
}

type request struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

type response struct {
	Message string `json:"message"`
	Role    string `json:"role"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func writeResponse(w http.ResponseWriter, res models.DataResponse, statusCode int) {
	b, _ := json.Marshal(res)
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
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
	if err != nil {
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
			return nil, fmt.Errorf(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
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

	switch acc := acc.(type) {
	case models.Merchant:
		claims = JWTData{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				Id:        acc.ID,
				IssuedAt:  time.Now().Unix(),
				Subject:   "merchant",
			},
			//map[string]string{
			//	"Id": merchant.ID,
			//	"Role": "merchant",
			//},
		}
	case models.User:
		claims = JWTData{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				Id:        acc.ID,
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

func EmailExists(w http.ResponseWriter, req *http.Request) {
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
	_, _ = w.Write(b)
}

func GenerateTokenByFacebookAccount(reqData models.FacebookTokenData) (interface{}, error) {
	user, err := services.GetFaceBookUserInfo(reqData)

	var response models.DataResponse

	if err != nil {
		log.Println(err)
		return nil, err
	}

	userInfo, err := services.GetUserByFbId(user.(models.FacebookUserInfo).ID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if userInfo == nil {
		//no account in user table
		userInfo, err = services.GetMerchantByFacebookID(user.(models.FacebookUserInfo).ID)

		if err != nil {
			return nil, err
		}

		if userInfo == nil {
			//no account in both user and merchant tables
			response.Success = true
			response.Type = "register"
			response.Message = reqData.AccessToken

			return response, nil
		} else {
			//merchant account exists
			token, err := createToken(userInfo)
			if err != nil {
				return nil, err
			}

			response.Success = true
			response.Type = "login"
			response.Role = c.Merchant
			response.Message = token

			return response, nil

		}
	} else {
		//user account exists
		token, err := createToken(userInfo)
		if err != nil {
			return nil, err
		}

		response.Success = true
		response.Type = "login"
		response.Role = c.User
		response.Message = token

		return response, nil
	}

}

func AuthenticateFacebook(w http.ResponseWriter, req *http.Request, userType string) {
	var reqData models.FacebookTokenData
	jsonEncoder := json.NewEncoder(w)

	err := json.NewDecoder(req.Body).Decode(&reqData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: "Incorrect data types"})
		return
	}

	switch userType {
	case "":
		response, err := GenerateTokenByFacebookAccount(reqData)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: err.Error()})
		} else {
			_ = jsonEncoder.Encode(response)
		}
	case c.Merchant:
		// do sign up merchant
		fmt.Printf("%+v", reqData)
		w.WriteHeader(http.StatusBadRequest)
		_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: "shit happened"})
	case c.User:
		//do sign up user
		fmt.Printf("%+v", reqData)
		w.WriteHeader(http.StatusBadRequest)
		_ = jsonEncoder.Encode(models.DataResponse{Success: false, Message: "shit happened"})
	}
}

func PhoneNumberExists(w http.ResponseWriter, req *http.Request) {
	var r request
	_ = json.NewDecoder(req.Body).Decode(&r)

	if merchant, _ := services.GetMerchantByPhoneNumber(r.PhoneNumber); merchant != nil {
		_ = json.NewEncoder(w).Encode(response{Message: "true", Role: c.Merchant})
	} else {
		if merchant, _ = services.GetUserByPhoneNumber(r.PhoneNumber); merchant != nil {
			_ = json.NewEncoder(w).Encode(response{Message: "true", Role: c.User})
		} else {
			_ = json.NewEncoder(w).Encode(response{Message: "false"})
		}
	}
}

func GetUserInfo(w http.ResponseWriter, req *http.Request) {
	token, err := getToken(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	claims, err := AuthenticateToken(token)

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
		_ = json.NewEncoder(w).Encode(res)

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
		_ = json.NewEncoder(w).Encode(res)
	}
}

type emailReq struct {
	Email string `json:"email"`
	Pin   string `json:"pin"`
}

type phoneReq struct {
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}

func SendRandomPinSMS(w http.ResponseWriter, req *http.Request) {
	var data phoneReq
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	pin := strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10))

	fmt.Println(pin)
	services.RedisClient.Set(data.PhoneNumber, pin, 5*time.Minute)
	services.SendSMSMessage(data.PhoneNumber, pin)
}

func VerifySMSPin(w http.ResponseWriter, req *http.Request) {
	var data phoneReq
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	if dbPin, err := services.RedisClient.Get(data.PhoneNumber).Result(); dbPin != data.Pin {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: "Invalid pin number"})
		return
	}

	services.RedisClient.Del(data.PhoneNumber)
	_ = json.NewEncoder(w).Encode(models.DataResponse{Success: true})
}

func SendRandomPinEmail(w http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)

	var data emailReq
	err := json.Unmarshal(b, &data)

	if err != nil {
		panic(err)
	}

	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	pin := strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10)) + strconv.Itoa(r1.Intn(10))

	fmt.Println(pin)

	services.RedisClient.Set(data.Email, pin, 5*time.Minute)
	_ = services.CreateEmailNotification(data.Email, "", pin)
}

func VerifyEmailPin(w http.ResponseWriter, req *http.Request) {
	b, _ := ioutil.ReadAll(req.Body)

	var data emailReq
	err := json.Unmarshal(b, &data)

	if err != nil {
		panic(err)
	}

	var res models.DataResponse
	if dbPin, err := services.RedisClient.Get(data.Email).Result(); dbPin != data.Pin {
		res.Success = false

		if err != nil {
			res.Message = err.Error()
		} else {
			res.Message = "invalid pin number"
		}

		w.WriteHeader(http.StatusBadRequest)
	} else {
		services.RedisClient.Del(data.Email)
		res.Success = true
	}

	_ = json.NewEncoder(w).Encode(res)
}
