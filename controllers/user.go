package controllers

import (
	auth "git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"net/http"
)

func UserHandler(w http.ResponseWriter, req *http.Request, objectID string) {
	if claims, err := auth.AuthenticateTokenMiddleWare(w, req); err != nil && claims.Role != "admin" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case "GET":
	case "POST":
	}
}
