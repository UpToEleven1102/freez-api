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

func handler(w http.ResponseWriter, req *http.Request) {
	repository, objectID := urlMatch(req.URL.Path)

	w.Header().Set("Content-type", "application/json")
	switch repository{
	case "merchants":
		controllers.MerchantHandler(w, req, objectID)
	default:
		http.NotFound(w, req)
	}
}

func main() {
	port:=getPort()

	http.HandleFunc("/api/", handler)
	fmt.Printf("Running on port %s \n", port)
	http.ListenAndServe(port, nil)
}

func getPort() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = ":8080"
	}
	return port
}
