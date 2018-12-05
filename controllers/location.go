package controllers

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func LocationHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) (err error){
	switch req.Method {
	case "POST":
		if objectID == "nearby" {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
			var location models.Location

			err = json.Unmarshal(body, &location.Location)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			merchants, err := services.GetNearMerchantsLastLocation(location)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

			b, _ := json.Marshal(merchants)
			w.Write(b)
		} else if len(objectID) == 0 {
			//api/location : used for merchant to push location to the db
			if claims.Role != config.Merchant {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return nil
			}

			body, _ := ioutil.ReadAll(req.Body)
			var location models.Location
			location.MerchantID = claims.Id

			err = json.Unmarshal(body, &location.Location)

			b, err := json.Marshal(location)
			fmt.Println(string(b))

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			err = services.AddNewLocation(location)
			if	err!=nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			http.NotFound(w, req)
		}

	case "GET":
		if len(objectID) == 0 {
			//get nearby merchants for user
		}
		// else get latest location for merchantID
		//TODO: include merchant info into response
		location, err := services.GetLastPositionByMerchantID(objectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		b, _ := json.Marshal(location)

		w.Write(b)
	}
	return nil
}