package main

import (
	"fmt"
	c "git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/controllers"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"github.com/tbalthazar/onesignal-go"
	"net/http"
	"os"
	"strings"
)

func init() {
	c.SetEnv()
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

	var err error
	switch repository {
	case c.Merchant:
		err = identity.AuthorizeMiddleware(w, req, objectID, controllers.MerchantHandler)
	case c.User:
		err = identity.AuthorizeMiddleware(w, req, objectID, controllers.UserHandler)
	case c.Request:
		err = identity.AuthorizeMiddleware(w, req, objectID, controllers.RequestHandler)
	case c.Location:
		err = identity.AuthorizeMiddleware(w, req, objectID, controllers.LocationHandler)
	default:
		http.NotFound(w, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
	port := getPort()
	_ = onesignal.NewClient(nil)

	http.HandleFunc("/api/", apiHandler)

	http.HandleFunc("/auth/", authHandler)

	fmt.Printf("Running on port %s \n", port)
	panic(http.ListenAndServe(port, nil))
}
