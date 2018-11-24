package services

import (
	"../models"
	"errors"
	"fmt"
	"encoding/json"
)

func CreateRequest(request models.Request) error {
	user, err := GetUserByEmail(request.Email)
	if err != nil || user == nil {
		return errors.New("bad request")
	}

	req, _ := json.Marshal(request)

	fmt.Println(string(req))
	point:= fmt.Sprintf(`POINT(%f %f)`, request.Lat, request.Long)

	_, err = DB.Exec(`INSERT INTO request (user_id, location) VALUES (?, ST_GeomFromText(?))`,user.(models.User).ID, point)

	return err
}

func GetRequest(userID string) (interface{}, error) {
	r, err := DB.Query(`SELECT user_id, ST_AsText(location) FROM request WHERE user_id=?`, userID)

	if err != nil {
		return nil, err
	}
	var id, point string
	if r.Next() {
		r.Scan(&id, &point)
		fmt.Println(id, point)
	}
	return nil, nil
}