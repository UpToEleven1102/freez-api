package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MerchantHandler(w http.ResponseWriter, req *http.Request, objectID string, claims models.JwtClaims) error {
	var response models.DataResponse
	fmt.Printf("IP address %s ", req.RemoteAddr)

	switch req.Method {
	case "GET":
		switch objectID {
		// @Title get merchants (Admin only)
		// @Description get the list of all merchants
		// @Accept  json
		// @Success 200 {array}  my_api.model.Merchant
		// @Resource /api
		// @Router /api/merchants [get]
		case "":
			if claims.Role != "admin" {
				return errors.New("Failed to authorize")
			}
			merchant, err := services.GetMerchantByEmail(objectID)
			if err != nil {
				return err
			}
			b, _ := json.Marshal(merchant)
			_, _ = w.Write(b)
			// @Title get presign-url
			// @Description get presign-url to upload image (jpg only)
			// @Accept  json
			// @Success 200 my_api.model.DataResponse
			// @Router /api/merchants/presign-url [get]
		case "presign-url":
			fileName := fmt.Sprint(claims.Id, "-profile.jpg")
			url, err := services.GeneratePreSignedUrl(fileName)

			var response models.DataResponse

			if err != nil {
				response.Success = false
				response.Message = err.Error()
			} else {
				response.Success = true
				response.Message = url
			}

			b, _ := json.Marshal(response)
			_, _ = w.Write(b)
		case "product":
			products, err := services.GetProducts(claims.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(products)
		case "notification":
			notifications, err := services.GetMerchantNotifications(claims.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(notifications)

		case "order":
			orders, err := services.GetOrderPaymentByMerchantId(claims.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:err.Error()})
				return nil
			}
			_ = json.NewEncoder(w).Encode(orders)

		default:
			objectID, param := getUrlParam(objectID)
			if param == "" {
				http.NotFound(w, req)
				return nil
			}

			switch objectID {
			case "notification":
				id, err := strconv.ParseInt(param, 0, 64)
				if err != nil {
					log.Println(err)
					http.NotFound(w, req)
					return nil
				}

				notification, err := services.GetMerchantNotificationById(id)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
					return nil
				}

				_ = json.NewEncoder(w).Encode(notification)
			}
		}

	case "POST":
		switch objectID {
		case "update-status":
			id := claims.Id
			err := services.ChangeOnlineStatus(id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return nil
			}
		case "product":
			var product models.Product

			err := json.NewDecoder(req.Body).Decode(&product)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}
			product.MerchantId = claims.Id
			err = services.CreateProduct(product)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

		case "food-type":
			var merchant models.Merchant

			err := json.NewDecoder(req.Body).Decode(&merchant)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()})
				return nil
			}
			merchant.ID = claims.Id
			err = services.UpdateFoodType(merchant)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message: err.Error()})
				return nil
			}

		case "product-presign-url":
			var product models.Product

			err := json.NewDecoder(req.Body).Decode(&product)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}
			var fileName string
			if product.Image != "" {
				arr := strings.Split(product.Image, "/")
				if arr[0] == "https:" {
					fileName = strings.Split(arr[len(arr)-1], "?")[0]
				}
			} else {
				fileName = fmt.Sprintf("%s-%d.jpg", claims.Id, time.Now().UnixNano())
			}

			url, err := services.GeneratePreSignedUrl(fileName)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			_ = json.NewEncoder(w).Encode(models.DataResponse{Success: true, Message: url})

		default:
			http.NotFound(w, req)
		}

	case "PUT":
		switch objectID {
		case "update-profile":

			b, err := ioutil.ReadAll(req.Body)

			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			var merchant models.Merchant
			err = json.Unmarshal(b, &merchant)
			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}
			merchant.ID = claims.Id

			err = services.UpdateMerchant(merchant)

			if err != nil {
				response.Success = false
				if strings.Contains(err.Error(), "Error 1062") {
					response.Message = "Email is currently in use!"
				} else {
					response.Message = err.Error()
				}

				sendResponse(w, response, http.StatusBadRequest)
			}

		case "order":
			var order models.OrderEntity
			err := json.NewDecoder(req.Body).Decode(&order)
			if err != nil {
				panic(err)
			}

			err = services.UpdateOrder(order)
			if err != nil {
				panic(err)
			}

		case "product":
			var product models.Product

			err := json.NewDecoder(req.Body).Decode(&product)

			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}
			product.MerchantId = claims.Id
			err = services.UpdateProduct(product)
			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

		case "notification":
			var notification models.MerchantNotification
			err := json.NewDecoder(req.Body).Decode(&notification)

			if err != nil {
				log.Println(err.Error())
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

			err = services.UpdateMerchantNotification(notification)

			if err != nil {
				log.Println(err.Error())
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
				return nil
			}

		default:
			http.NotFound(w, req)
		}

	case "DELETE":
		switch objectID {
		case "product":
			var data models.Product

			_ = json.NewDecoder(req.Body).Decode(&data)

			err := services.DeleteProduct(data)

			if err != nil {
				_ = json.NewEncoder(w).Encode(models.DataResponse{Success: false, Message: err.Error()})
			}
		}

	default:
		http.NotFound(w, req)
	}

	return nil
}
