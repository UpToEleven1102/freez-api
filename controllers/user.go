package controllers

import (
	"encoding/json"
	"errors"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func UserHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	if claims.Role != config.User {
		return errors.New("Unauthenticated")
	}

	switch req.Method {
	case "GET":
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
	}

	return nil
}
