package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func sendResponse(w http.ResponseWriter, response models.DataResponse, status int) {
	w.WriteHeader(status)
	b, _ := json.Marshal(response)
	w.Write(b)
}

func getUrlParam(objectID string) (objectId string, param string) {
	arr := strings.Split(objectID, "/")
	if len(arr) > 1 {
		return arr[0], arr[1]
	}
	return objectID, ""
}

func UserHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	if claims.Role != config.User {
		return errors.New("Unauthenticated")
	}

	var response models.DataResponse

	switch req.Method {
	case "GET":
		switch objectID {
		case "presign-url":
			fileName := fmt.Sprint(claims.Id, "-profile.jpg")
			url, err := services.GeneratePreSignedUrl(fileName)

			if err != nil {
				response.Success = false
				response.Message = err.Error()
				sendResponse(w, response, http.StatusInternalServerError)
				return nil
			} else {
				response.Success = true
				response.Message = url
			}

			b, _ := json.Marshal(response)
			_, _ = w.Write(b)
		case "favorite":
			r, err := services.GetFavorites(claims.Id)

			if err != nil {
				response.Success = false
				response.Message = err.Error()
				sendResponse(w, response, http.StatusInternalServerError)
				return nil
			}

			b, _ := json.Marshal(r)
			w.Write(b)
		case "notification":
			notifications, err := services.GetUserNotifications(claims.Id)
			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(notifications)
		default:
			objectID, param := getUrlParam(objectID)
			if param == "" {
				http.NotFound(w, req)
				return nil
			}

			switch objectID {
			case "notification":
				id, err := strconv.ParseInt(param, 0, 64)

				if err != nil {
					log.Println(err)
					http.NotFound(w, req)
					return nil
				}
				notification, err := services.GetUserNotificationById(id)

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false})
				}

				_ = json.NewEncoder(w).Encode(notification)
			}
		}

	case "POST":
		switch objectID {
		case "location":
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
			var location models.LongLat
			_ = json.Unmarshal(b, &location)

			var user models.User
			user.ID = claims.Id
			user.LastLocation = location
			_, err = services.UpdateUserLocation(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

		case "favorite":
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			var data models.RequestData

			_ = json.Unmarshal(b, &data)
			data.UserId = claims.Id

			err = services.AddFavorite(data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "charge":
			type chargeData struct {
				StripeToken string  `json:"stripeToken"`
				Amount      float64 `json:"amount"`
			}
			var data chargeData

			_ = json.NewDecoder(req.Body).Decode(&data)

			res, err := services.StripeCharge(data.StripeToken, "Testing", data.Amount)

			if err != nil {
				log.Println(err)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			fmt.Printf("%+v", res)
			_ = json.NewEncoder(w).Encode(models.DataResponse{Success:true})

		default:
			http.NotFound(w, req)
		}

	case "DELETE":
		switch objectID {
		case "favorite":
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			var data models.RequestData

			_ = json.Unmarshal(b, &data)
			data.UserId = claims.Id

			err = services.RemoveFavorite(data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

	case "PUT":
		switch objectID {
		case "update-profile":
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			var user models.User

			err = json.Unmarshal(b, &user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
			user.ID = claims.Id

			err = services.UpdateUser(user)

			fmt.Println(err)
			if err != nil {
				response.Success = false
				if strings.Contains(err.Error(), "Error 1062") {
					response.Message = "Email is currently in use!"
				} else {
					response.Message = err.Error()
				}

				sendResponse(w, response, http.StatusBadRequest)
			}
		case "notification":
			var notification models.UserNotification

			err := json.NewDecoder(req.Body).Decode(&notification)

			fmt.Printf("%+v\n", notification)

			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			err = services.UpdateUserNotification(notification)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}
		}
	}

	return nil
}
