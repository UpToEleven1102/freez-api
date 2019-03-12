// @APIVersion 1.0.0
// @APITitle My Cool Freez API
// @APIDescription My API usually works as expected (sometimes it doesn't).
// @BasePath http://35.162.158.187/
package main

import (
	"fmt"
	c "git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/controllers"
	"git.nextgencode.io/huyen.vu/freez-app-rest/identity"
	"github.com/joho/godotenv"
	"github.com/tbalthazar/onesignal-go"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func urlMatch(url string) (repository string, objectID string) {
	fragments := strings.SplitN(url, "/", -1)
	repository = fragments[2]
	objectID = ""
	if len(fragments) == 4 {
		objectID = fragments[3]
	} else if len(fragments) > 4 {

		objectID = fragments[3]
		for i := 4; i < len(fragments); i++ {
			objectID += "/" + fragments[i]
		}
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
	case c.Stripe:
		err = identity.AuthorizeMiddleware(w, req, objectID, controllers.StripeOpsHandler)

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

	// @SubApi Main API [/api]
	http.HandleFunc("/api/", apiHandler)

	// @SubApi Auth API [/auth]
	http.HandleFunc("/auth/", authHandler)

	// @SubApi Websocket routes for users [/socker/user]
	http.Handle("/socket/user", websocket.Handler(controllers.UserWebSocketHandler))

	// @SubApi Websocket routes for merchants [/socket/merchant]
	http.Handle("/socket/merchant", websocket.Handler(controllers.MerchantWebSocketHandler))

	fmt.Printf("Running on port %s \n", port)
	panic(http.ListenAndServe(port, nil))
}
