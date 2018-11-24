package controllers

import (
	"net/http"
	"io/ioutil"
	"../models"
	"encoding/json"
	"../services"
	"../identity"
)

func RequestHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	_ , err :=identity.AuthenticateTokenMiddleWare(w, req)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "POST":
		if len(objectID) > 0 {
			http.NotFound(w, req)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var request models.Request
		json.Unmarshal(body, &request)

		err = services.CreateRequest(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	case "GET":
		r, err := services.GetUserByEmail(objectID)
		if err != nil || r == nil {
			http.Error(w, "user not exists", http.StatusBadRequest)
			return
		}
		user := r.(models.User)
		r, err = services.GetRequest(user.ID)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		request := r.(models.Request)

		b, _ := json.Marshal(request)
		w.Write(b)
	}

}