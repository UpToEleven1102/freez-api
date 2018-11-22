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

func urlMatch(url string) (repository string, objectID string) {
	fragments := strings.SplitN(url, "/", -1)
	repository = fragments[2]
	objectID = ""
	 if len(fragments) > 3 {
		objectID = fragments[3]
	}

	return repository, objectID
}

func getPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = ":8080"
	}
	return port
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	repository, objectID := urlMatch(req.URL.Path)

	w.Header().Set("Content-type", "application/json")
	switch repository{
	case "merchants":
		controllers.MerchantHandler(w, req, objectID)
	default:
		http.NotFound(w, req)
	}
}

func authHandler(w http.ResponseWriter, req *http.Request) {
	route, userType := urlMatch(req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	if req.Method != "POST" {
		http.NotFound(w, req)
	}

	controllers.AuthHandler(w, req, route, userType)
}

func main() {
	port:=getPort()

	http.HandleFunc("/api/", apiHandler)

	http.HandleFunc("/auth/", authHandler)

	fmt.Printf("Running on port %s \n", port)
	http.ListenAndServe(port, nil)
}
