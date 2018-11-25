package controllers

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"net/http"
)

func LocationHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) (err error){
	switch req.Method {
	case "POST":
		if claims.Role != config.Merchant {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}

		if len(objectID) > 0 {
			http.NotFound(w, req)
			return nil
		}

		body, _ := ioutil.ReadAll(req.Body)
		var location models.Location
		location.MerchantID = claims.Id

		err = json.Unmarshal(body, &location.Location)

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

	case "GET":
		if len(objectID) == 0 {
			//get near by merchants for user
		}
		// else get latest location for merchantID

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