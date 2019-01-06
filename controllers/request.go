package controllers

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/services"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func RequestHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	switch req.Method {
	case "POST":
		if claims.Role != "user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}

		if len(objectID) > 0 {
			http.NotFound(w, req)
			return nil
		}

		body, _ := ioutil.ReadAll(req.Body)

		var request models.Request
		err := json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		request.UserId = claims.Id

		err = services.CreateRequest(request, claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}

		w.WriteHeader(http.StatusOK)
	case "GET":
		switch objectID {
		case "":
			if claims.Role == config.User {
				r, err := services.GetRequestedMerchantByUserID(claims.Id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}
				b, _ := json.Marshal(r)
				_, _ = w.Write(b)

			} else if claims.Role == config.Merchant {
				r, err := services.GetRequestInfoByMerchantId(claims.Id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return nil
				}

				b, _ := json.Marshal(r)
				_, _ = w.Write(b)
			}

			return nil
		case "user":
			id := claims.Id
			r, err := services.GetRequestByUserID(id)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return nil
			}

			if r == nil {
				w.Write(nil)
				return nil
			}
			request := r.(models.Request)

			b, _ := json.Marshal(request)
			w.Write(b)
		case "merchant":
			id := claims.Id
			r, err := services.GetRequestByMerchantID(id)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return nil
			}

			if r == nil {
				w.Write(nil)
				return nil
			}
			request := r.(models.Request)

			b, _ := json.Marshal(request)
			w.Write(b)

		default:
			param, err := strconv.ParseInt(objectID, 0, 64)

			if err != nil {
				log.Println(err)
				http.NotFound(w, req)
				return nil
			}

			if	claims.Role == config.Merchant {
				request, err := services.GetRequestInfoById(int(param), claims.Id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false})
					return nil
				}

				_ = json.NewEncoder(w).Encode(request)
			}
		}
	case "PUT":
		if claims.Role == config.Merchant {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}

			var request models.RequestEntity
			json.Unmarshal(b, &request)

			fmt.Printf("%+v", request)

			if request.Accepted != 0 && request.Accepted != 1 {
				http.Error(w, "accepted param must be 0 or 1", http.StatusBadRequest)
				return nil
			}

			request.MerchantID = claims.Id
			err = services.UpdateRequestAccepted(request, claims)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
		}
	case "DELETE":
		if len(objectID) > 0 {
			http.NotFound(w, req)
			return nil
		}

		err := services.RemoveRequestsByUserID(claims.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
	}

	return nil
}
