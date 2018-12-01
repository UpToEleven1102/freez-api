package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func RequestHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	switch req.Method {
	case "POST":
		if claims.Role != "user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}

		if len(objectID) > 0 {
			http.NotFound(w, req)
			return nil
		}

		body, _ := ioutil.ReadAll(req.Body)

		var request models.Request
		err := json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		request.UserId = claims.Id

		err = services.CreateRequest(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		w.WriteHeader(http.StatusOK)
	case "GET":
		switch objectID {
		case "" :
			r, err := services.GetRequests()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

			b, _ := json.Marshal(r)
			w.Write(b)
			return nil
		case "user":
			id := claims.Id
			r, err := services.GetRequestByUserID(id)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return nil
			}

			if r == nil {
				w.Write(nil)
				return nil
			}
			request := r.(models.Request)

			b, _ := json.Marshal(request)
			w.Write(b)
		case "merchant":
		default:
			http.NotFound(w,req)
		}
	case "DELETE":
		if len(objectID) > 0 {
			http.NotFound(w, req)
			return nil
		}

		err := services.RemoveRequestsByUserID(claims.Id)
		if	err!=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
	}

	return nil
}
