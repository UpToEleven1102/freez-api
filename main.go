package main

import (
	"./controllers"
	"./config"
	"net/http"
	"os"
	"strings"

	"fmt"
)

func init() {
	config.ConfigEnv()
}

func urlMatch(url string) (repository string, attribute string, objectID string) {
	fragments := strings.SplitN(url, "/", -1)
	repository = fragments[2]
	objectID = ""
	attribute = ""
	if len(fragments) == 4 {
		objectID = fragments[3]
	} else if len(fragments) > 4 {
		attribute = fragments[3]
		objectID = fragments[4]
	}

	return repository, attribute, objectID
}

func getPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = ":8080"
	}
	return port
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	repository, attribute,  objectID := urlMatch(req.URL.Path)

	w.Header().Set("Content-type", "application/json")
	switch repository{
	case "merchants":
		controllers.MerchantHandler(w, req, attribute, objectID)
	default:
		http.NotFound(w, req)
	}
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	route, _, _ := urlMatch(req.URL.Path)
	w.Header().Set("Content-type", "application/json")
	if req.Method != "POST" {
		http.NotFound(w, req)
	}

	switch route {
	case "signup":
		controllers.SignUp(w, req)
	case "signin":
		controllers.SignIn(w, req)
	default:
		http.NotFound(w, req)
	}
}

func main() {
	port:=getPort()

	http.HandleFunc("/api/", apiHandler)

	http.HandleFunc("/auth/", authHandler)

	fmt.Printf("Running on port %s \n", port)
	http.ListenAndServe(port, nil)
}
