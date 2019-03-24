package identity

import (
	"encoding/json"
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

func SignInMerchant(w http.ResponseWriter, req *http.Request) {
	var credentials Credentials
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := services.GetMerchantByEmail(credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if res != nil {
		merchant := res.(models.Merchant)
		err = bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(credentials.Password))
		if err == nil {
			token, err := createToken(merchant)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(models.DataResponse{Success:true, Message:token})
			return
		}
	}
	_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:"Credentials Invalid"})
}

func SignUpMerchantFB(reqData models.FacebookTokenData) (response models.DataResponse, err error) {
	response.Success = false
	response.Type = config.SignUp
	response.Role = config.Merchant

	fbInfo, err := services.GetFaceBookUserInfo(reqData)

	userInfo := fbInfo.(models.FacebookUserInfo)
	if err != nil {
		response.Message = err.Error()
		return response, err
	}

	merchant := models.Merchant{
		Email: userInfo.Email,
		Image: userInfo.Picture,
		Name: userInfo.Name,
		PhoneNumber: reqData.PhoneNumber,
		FacebookID: userInfo.ID,
		Password: reqData.Password,
		Mobile: reqData.IsMobile,
		StripeID: reqData.Stripe_Token,
		Category: reqData.Category,
	}

	acc, err := services.StripeConnectCreateAccount(merchant)
	if err != nil {
		response.Message = err.Error()
		return response, err
	}
	merchant.StripeID = acc.ID

	merchant, err = services.CreateMerchant(merchant)

	fmt.Printf("merchant: %+v \n", merchant)

	if err != nil {
		fmt.Println(err.Error())
		response.Message = err.Error()
		return response, err
	}

	token, err := createToken(merchant)
	if err != nil {
		response.Message = err.Error()
		return response, err
	}

	response.Success = true
	response.Message = token

	return response, err
}

func SignUpMerchant(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	var response models.DataResponse
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusBadRequest)
		fmt.Println(response)
		return
	}

	var merchant models.Merchant

	err = json.Unmarshal(body, &merchant)

	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusBadRequest)
		return
	}

	acc, err := services.StripeConnectCreateAccount(merchant)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:err.Error()})
		return
	}

	//res, err := services.StripeCharge(merchant.StripeID, "Application Fee - Freeze App", 5)

	//if err != nil {
	//	log.Println(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:err.Error()})
	//	return
	//}

	//log.Println(res)

	merchant.StripeID = acc.ID

	merchant, err = services.CreateMerchant(merchant)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		writeResponse(w, response, http.StatusInternalServerError)
		fmt.Println(response)
		return
	}

	token, _ := createToken(merchant)
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
		Message string `json:"message"`
	}{true, token})
}
