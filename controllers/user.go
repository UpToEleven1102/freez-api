package controllers

import "net/http"

func UserHandler(w http.ResponseWriter, req *http.Request, objectID string){
	method:= req.Method

	switch method {
	case "GET":
	case "POST":
	}
}