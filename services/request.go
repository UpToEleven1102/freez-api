package services

import (
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"strconv"
	"strings"
)

func CreateRequest(request models.Request) error {
	//use userId from claims instead

	point := fmt.Sprintf(`POINT(%f %f)`, request.Location.Lat, request.Location.Long)

	if r, err := GetMerchantById(request.MerchantID); err != nil || r == nil {
		return errors.New("merchant doesn't exist")
	}

	_, err := DB.Exec(`INSERT INTO request (user_id, merchant_id, location) VALUES (?, ?, ST_GeomFromText(?))`, request.UserId, request.MerchantID, point)
	return err
}

func getLatLong(point string) (lat float32, long float32, err error) {
	ptArr := strings.Split(strings.Replace(point, ")", "", -1), "(")
	if len(ptArr) < 2 {
		return 0, 0, errors.New("index out of range")
	}
	ptArr = strings.Split(ptArr[1], " ")
	if len(ptArr) < 2 {
		return 0, 0, errors.New("index out of range")
	}
	lat64, _ := strconv.ParseFloat(ptArr[0], 32)
	long64, _ := strconv.ParseFloat(ptArr[1], 32)
	lat = float32(lat64)
	long = float32(long64)
	return lat,long, nil
}

func GetRequestByUserID(userID string) (interface{}, error) {
	//r, err := DB.Query(`SELECT id, ST_AsText(location) FROM request WHERE user_id=?`, userID)
	//
	//if err != nil {
	//	return nil, err
	//}
	//var id int
	//var point string
	//if r.Next() {
	//	err = r.Scan(&id, &point)
	//
	//	user, err := GetUserById(userID)
	//	if err != nil || user == nil{
	//		return nil, err
	//	}
	//
	//	var request models.Request
	//	request.Email = user.(models.User).Email
	//	request.Location.Lat, request.Location.Long, err = getLatLong(point)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return request, nil
	//}
	return nil, nil
}

func GetRequests() (requests []interface{}, err error) {
	_, err = DB.Query(`SELECT email, ST_AsText(location) FROM request r JOIN user u ON r.user_id=u.id;`)
	if err != nil {
		return nil, err
	}

	//var location, email string
	//
	//for r.Next()  {
	//
	//	err = r.Scan(&email, &location)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	var request models.Request
	//
	//	request.Email = email
	//	request.Location.Lat, request.Location.Long, err = getLatLong(location)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	requests = append(requests, request)
	//}
	return requests, nil
}

func RemoveRequestsByUserID(userID string) (err error) {
	_, err = DB.Exec(`DELETE FROM request WHERE user_id=?`, userID)
	return err
}