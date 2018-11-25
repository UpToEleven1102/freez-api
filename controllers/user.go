package controllers

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"net/http"
)

func UserHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error{
 	switch req.Method {
	case "GET":
	case "POST":
	}

 	return nil
}
