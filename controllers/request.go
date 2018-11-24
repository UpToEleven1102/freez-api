package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func RequestHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	claims, err := identity.AuthenticateTokenMiddleWare(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "POST":
		if len(objectID) > 0 {
			http.NotFound(w, req)
			return
		}

		body, _ := ioutil.ReadAll(req.Body)
		var request models.Request
		err := json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = services.CreateRequest(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	case "GET":
		switch len(objectID) {
		case 0 :
			r, err := services.GetRequests()
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			b, _ := json.Marshal(r)
			w.Write(b)
		default:
			r, err := services.GetUserByEmail(objectID)
			if err != nil || r == nil {
				http.Error(w, "user not exists", http.StatusBadRequest)
				return
			}
			user := r.(models.User)
			r, err = services.GetRequestByUserID(user.ID)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			request := r.(models.Request)

			b, _ := json.Marshal(request)
			w.Write(b)
		}
	case "DELETE":
		if len(objectID) > 0 {
			http.NotFound(w, req)
			return
		}

		err = services.RemoveRequestsByUserID(claims.Id)
		if	err!=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
