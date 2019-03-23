package identity

import (
	"encoding/json"
	"git.nextgencode.io/huyen.vu/freez-app-rest/config"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"git.nextgencode.io/huyen.vu/freez-app-rest/services"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"fmt"
)

func SignUpUserFB(reqData models.FacebookTokenData) (response models.DataResponse, err error) {
	response.Success = false
	response.Role = config.User
	response.Type = config.SignUp

	fbInfo, err := services.GetFaceBookUserInfo(reqData)

	if err != nil {
		response.Message = err.Error()
		return response, err
	}

	userInfo := fbInfo.(models.FacebookUserInfo)

	user, err := services.CreateUserFB(models.User{
		Image: userInfo.Picture,
		Name: userInfo.Name,
		Email:userInfo.Email,
		PhoneNumber: reqData.PhoneNumber,
		FacebookID:userInfo.ID,
		Password:reqData.Password,
	})

	if err != nil {
		fmt.Printf("%s \n", err.Error())
		response.Message = err.Error()
		return response, err
	}

	token, err := createToken(user)
	if err != nil {
		fmt.Printf("%s \n", err.Error())
		response.Message = err.Error()
		return response, err
	}

	response.Success = true
	response.Message = token

	return response, err
}

func SignUpUser(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var user models.User
	_ = json.Unmarshal(body, &user)

	fmt.Printf("%+v", user)

	r, err := services.CreateUser(user)
	if err != nil {
		http.Error(w, "account exists", http.StatusBadRequest)
		return
	}

	user = r.(models.User)
	token, err := createToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(token)
	w.WriteHeader(http.StatusCreated)
	_ , _ = w.Write(b)
}

func SignInUser(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var credentials Credentials
	_ = json.Unmarshal(body, &credentials)

	r, err := services.GetUserByEmail(credentials.Email)

	if err != nil || r == nil {
		_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:"Credentials Invalid"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(r.(models.User).Password), []byte(credentials.Password))

	if err != nil {
		_ = json.NewEncoder(w).Encode(models.DataResponse{Success:false, Message:"Credentials Invalid"})
		return
	}

	token, _ := createToken(r)
	_ = json.NewEncoder(w).Encode(models.DataResponse{Success:true, Message:token})
}
