package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"io/ioutil"
	"log"
	"net/http"
)

/*LocationHandler - HandleFunc for location route*/
func LocationHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) (err error){
	jsonEncoder := json.NewEncoder(w)
	switch req.Method {
	case "POST":
		switch objectID {
		case ""://api/location : used for merchant to push location to the db
			if claims.Role != config.Merchant {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return nil
			}

			jsonEncoder := json.NewEncoder(w)

			body, _ := ioutil.ReadAll(req.Body)
			var location models.Location
			location.Id = claims.Id

			err = json.Unmarshal(body, &location.Location)
			if err != nil {
				_ = jsonEncoder.Encode(models.DataResponse{Type:"error", Message:err.Error()})
				return nil
			}

			err = services.AddNewLocation(location)
			if err != nil {
				_ = jsonEncoder.Encode(models.DataResponse{Type:"error", Message:err.Error()})
				return nil
			}

			w.WriteHeader(http.StatusCreated)

		case "nearby":
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

			location.Id = claims.Id

			merchants, err := services.GetNearbyMerchantsLastLocation(location)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

			panic(jsonEncoder.Encode(merchants))

		default:
			objectID, param := getUrlParam(objectID)

			if len(param) == 0 {
				http.NotFound(w, req)
				return nil
			}

			switch objectID {
			case "nearby":
				var location models.Location
				err := json.NewDecoder(req.Body).Decode(&location)

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					return nil
				}

				location.Id = claims.Id

				merchant, err := services.GetMerchantInfoById(param, location)

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return nil
				}

				_ = json.NewEncoder(w).Encode(merchant)
			}
		}

	case "GET":
		//if len(objectID) == 0 {
		//	//get nearby merchants for user
		//}
		// else get latest location for merchantID
		//TODO: include merchant info into response
		location, err := services.GetLastPositionByMerchantID(objectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		panic(jsonEncoder.Encode(location))
	}
	return nil
}