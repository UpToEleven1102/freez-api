package services

import (
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"strconv"
	"strings"
)

func CreateRequest(request models.Request) error {
	user, err := GetUserByEmail(request.Email)
	if err != nil || user == nil {
		return errors.New("bad request")
	}

	point := fmt.Sprintf(`POINT(%f %f)`, request.Lat, request.Long)

	_, err = DB.Exec(`INSERT INTO request (user_id, location) VALUES (?, ST_GeomFromText(?))`, user.(models.User).ID, point)

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
	r, err := DB.Query(`SELECT id, ST_AsText(location) FROM request WHERE user_id=?`, userID)

	if err != nil {
		return nil, err
	}
	var id int
	var point string
	if r.Next() {
		err = r.Scan(&id, &point)

		user, err := GetUserById(userID)
		if err != nil || user == nil{
			return nil, err
		}

		var request models.Request
		request.Email = user.(models.User).Email
		request.Lat, request.Long, _ = getLatLong(point)
		return request, nil
	}
	return nil, nil
}

func GetRequests() (requests []interface{}, err error) {
	r, err := DB.Query(`SELECT id, user_id, ST_AsText(location) FROM request;`)
	var userID, location string
	var id int
	for r.Next()  {


		err = r.Scan(&id, &userID, &location)
		if err != nil {
			return nil, err
		}
		var request models.Request

		user, err := GetUserById(userID)
		if err != nil {
			return nil, err
		}

		if user != nil {
			request.Email = user.(models.User).Email
			request.Lat, request.Long, _ = getLatLong(location)
			requests = append(requests, request)
		}
	}
	return requests, err
}

func RemoveRequestsByUserID(userID string) (err error) {
	_, err = DB.Exec(`DELETE FROM request WHERE user_id=?`, userID)
	return err
}