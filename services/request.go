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

	point := fmt.Sprintf(`POINT(%f %f)`, request.Location.Long, request.Location.Lat)

	if r, err := GetMerchantById(request.MerchantID); err != nil || r == nil {
		return errors.New("merchant doesn't exist")
	}

	_, err := DB.Exec(`INSERT INTO request (user_id, merchant_id, location) VALUES (?, ?, ST_GeomFromText(?))`, request.UserId, request.MerchantID, point)
	return err
}

func getLongLat(point string) (long float32, lat float32, err error) {
	ptArr := strings.Split(strings.Replace(point, ")", "", -1), "(")
	if len(ptArr) < 2 {
		return 0, 0, errors.New("index out of range")
	}
	ptArr = strings.Split(ptArr[1], " ")
	if len(ptArr) < 2 {
		return 0, 0, errors.New("index out of range")
	}
	long64, _ := strconv.ParseFloat(ptArr[0], 32)
	lat64, _ := strconv.ParseFloat(ptArr[1], 32)
	lat = float32(lat64)
	long = float32(long64)
	return long, lat,nil
}

func GetRequestByUserID(userID string) (interface{}, error) {
	r, err := DB.Query(`SELECT user_id, merchant_id, ST_AsText(location) FROM request WHERE user_id=?`, userID)

	if err != nil {
		return nil, err
	}

	var point string
	if r.Next() {
		var request models.Request

		err = r.Scan(&request.UserId, &request.MerchantID, &point)
		request.Location.Long, request.Location.Lat, err = getLongLat(point)
		if err != nil {
			return nil, err
		}


		if err != nil {
			return nil, err
		}
		return request, nil
	}
	return nil, nil
}

func GetRequestByMerchantID(merchantID string) (interface{}, error) {
	r, err := DB.Query(`SELECT user_id, merchant_id, ST_ASTEXT(location) FROM request WHERE merchant_id=?`, merchantID)

	if err != nil {
		return nil, err
	}

	var point string

	if r.Next() {
		var request models.Request

		err = r.Scan(&request.UserId, &request.MerchantID, &point)
		request.Location.Long, request.Location.Lat, err = getLongLat(point)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}
		return request, nil
	}

	return nil, nil
}

func GetRequests() (interface{}, error){
	r, err := DB.Query(`SELECT user_id, merchant_id, ST_ASTEXT(location) FROM request`)

	if err != nil {
		return nil, err
	}

	var requests []models.Request
	var request models.Request
	var point string
	if r.Next() {
		err = r.Scan(&request.UserId, &request.MerchantID, &point)
		if err != nil {
			return nil, err
		}

		request.Location.Long, request.Location.Lat, _ = getLongLat(point)
		requests = append(requests, request)
	}

	return requests, nil
}

func RemoveRequestsByUserID(userID string) (err error) {
	_, err = DB.Exec(`DELETE FROM request WHERE user_id=?`, userID)
	return err
}