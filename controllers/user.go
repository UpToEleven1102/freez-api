package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func sendResponse(w http.ResponseWriter, response models.DataResponse, status int) {
	w.WriteHeader(status)
	b, _ := json.Marshal(response)
	w.Write(b)
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
			_,_ = w.Write(b)
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
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

		}
	}

	return nil
}
