package services

import (
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"log"
	"strconv"
	"strings"
)

const (
	notificationRequestMessage         = "New Request From User"
	notificationRequestTitle           = "New Request"
	notificationRequestDeclinedMessage = "Merchant is busy. Please request a bit later."
	notificationRequestAcceptedMessage = "Merchant is on the way"
)

func CreateRequest(request models.Request, claims models.JwtClaims) error {
	//use userId from claims instead

	point := fmt.Sprintf(`POINT(%f %f)`, request.Location.Long, request.Location.Lat)

	if r, err := GetMerchantById(request.MerchantID); err != nil || r == nil {
		return errors.New("merchant doesn't exist")
	}

	r, err := DB.Exec(`INSERT INTO request (user_id, merchant_id, location) VALUES (?, ?, ST_GeomFromText(?))`, request.UserId, request.MerchantID, point)
	rID, _ := r.LastInsertId()

	data := models.RequestData{UserId: "id---", Data: "S3cr3t"}
	if err == nil {
		_, err := CreateNotificationByUserId(request.MerchantID, notificationRequestTitle, notificationRequestMessage, claims, data)
		if err != nil {
			panic(err)
		}
		err = InsertMerchantNotification(request.MerchantID, 1, rID, notificationRequestMessage)
		if err != nil {
			panic(err)
		}
	}
	return err
}

func getLongLat(point string) (long float32, lat float32, err error) {
	if point == "" {
		return 0, 0, nil
	}
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
	return long, lat, nil
}

func GetRequestByUserID(userID string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, user_id, merchant_id, ST_AsText(location), comment, accepted, active FROM request WHERE user_id=?`, userID)
	defer r.Close()

	if err != nil {
		return nil, err
	}

	var point string
	if r.Next() {
		var request models.RequestEntity

		err = r.Scan(&request.ID ,&request.UserID, &request.MerchantID, &point, &request.Comment, &request.Accepted, &request.Active)
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

func GetRequestById(id int) (interface{}, error) {
	r, err := DB.Query(`SELECT id, user_id, merchant_id, ST_AsText(location), comment, accepted, active FROM request WHERE id=?`, id)
	defer r.Close()

	if err != nil {
		return nil, err
	}

	var point string
	if r.Next() {
		var request models.RequestEntity

		err = r.Scan(&request.ID ,&request.UserID, &request.MerchantID, &point, &request.Comment, &request.Accepted, &request.Active)
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

func GetRequestNotificationById(id int) (interface{}, error) {
	r, err := DB.Query(`SELECT r.id, r.comment, r.active, r.accepted, 
       								user_id, u.name, u.image, u.phone_number, u.email,
       								merchant_id, m.name, m.image, m.phone_number, m.email, m.online, m.mobile , ST_AsText(m.last_location), ST_DISTANCE_SPHERE(m.last_location, u.last_location) as distance   								
								FROM request r 
								JOIN user u 
									ON r.user_id=u.id
								JOIN merchant m
									ON r.merchant_id=m.id
								WHERE r.id=?`, id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer r.Close()

	var reqNotif models.RequestNotification
	var location string

	for r.Next() {
		err = r.Scan(&reqNotif.ID, &reqNotif.Comment, &reqNotif.Active, &reqNotif.Accepted,
			&reqNotif.User.ID, &reqNotif.User.Name, &reqNotif.User.Image, &reqNotif.User.PhoneNumber, &reqNotif.User.Email,
			&reqNotif.Merchant.MerchantID, &reqNotif.Merchant.Name, &reqNotif.Merchant.Image, &reqNotif.Merchant.PhoneNumber, &reqNotif.Merchant.Email, &reqNotif.Merchant.Online, &reqNotif.Merchant.Mobile, &location, &reqNotif.Merchant.Distance)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		reqNotif.Merchant.Accepted = reqNotif.Accepted
		reqNotif.Merchant.IsFavorite, err = isFavorite(models.RequestData{UserId:reqNotif.User.ID, Data: reqNotif.Merchant.MerchantID})

		if err != nil {
			log.Println(err)
		}

		return reqNotif, nil
	}

	return nil, nil
}

func GetRequestByMerchantID(merchantID string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, user_id, merchant_id, ST_AsText(location), comment, accepted, active FROM request WHERE merchant_id=?`, merchantID)
	defer r.Close()

	if err != nil {
		return nil, err
	}

	var point string

	if r.Next() {
		var request models.RequestEntity

		err = r.Scan(&request.ID ,&request.UserID, &request.MerchantID, &point, &request.Comment, &request.Accepted, &request.Active)
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

func GetRequestInfoById(id int, merchantId string) (interface{}, error) {
	position, _ := GetLastPositionByMerchantID(merchantId)
	var location string
	if position != nil {
		location = fmt.Sprintf(`POINT(%f %f)`, position.(models.Location).Location.Long, position.(models.Location).Location.Lat)
	}

	r, err := DB.Query(`SELECT r.id, user_id, merchant_id, name, email, phone_number, image, ST_ASTEXT(location), ST_DISTANCE_SPHERE(location, ST_GeomFromText(?)) as distance, comment, accepted, active
								FROM request r 
								  LEFT OUTER JOIN user u 
								    ON r.user_id=u.id 
								WHERE r.id=?`, location, id)
	defer r.Close()

	if err != nil {
		panic(err)
		return nil, err
	}

	location = ""

	if r.Next() {
		var request models.RequestInfo
		err = r.Scan(&request.ID, &request.UserId,&request.MerchantId, &request.Name, &request.Email, &request.PhoneNumber, &request.Image, &location, &request.Distance, &request.Comment, &request.Accepted, &request.Active)
		request.Location.Long, request.Location.Lat, _ = getLongLat(location)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		return request, nil
	}
	return nil, nil
}

func GetRequestInfoByMerchantId(merchantId string) (interface{}, error) {
	position, _ := GetLastPositionByMerchantID(merchantId)
	var location string
	if position != nil {
		location = fmt.Sprintf(`POINT(%f %f)`, position.(models.Location).Location.Long, position.(models.Location).Location.Lat)
	}

	r, err := DB.Query(`SELECT r.id, user_id, merchant_id, name, email, phone_number, image, ST_ASTEXT(location), ST_DISTANCE_SPHERE(location, ST_GeomFromText(?)) as distance, comment, accepted, active
								FROM request r 
								  LEFT OUTER JOIN user u 
								    ON r.user_id=u.id 
								WHERE merchant_id=? AND active=TRUE`, location, merchantId)
	defer r.Close()

	if err != nil {
		panic(err)
		return nil, err
	}

	var requests []models.RequestInfo
	location = ""

	for r.Next() {
		var request models.RequestInfo
		_ = r.Scan(&request.ID, &request.UserId,&request.MerchantId, &request.Name, &request.Email, &request.PhoneNumber, &request.Image, &location, &request.Distance, &request.Comment, &request.Accepted, &request.Active)
		request.Location.Long, request.Location.Lat, _ = getLongLat(location)
		requests = append(requests, request)
	}
	return requests, nil
}

func GetRequestedMerchantByUserID(userId string) (interface{}, error) {
	r, err := DB.Query(`SELECT online, m.email, m.name, mobile, m.phone_number, m.image, l.merchant_id, ST_AsText(l.location) as location, ST_DISTANCE_SPHERE(l.location, u.last_location) as distance, accepted
								FROM location l INNER JOIN (
								    SELECT merchant_id, MAX(ts) AS ts FROM location GROUP BY merchant_id
								  ) latest
								  ON l.ts=latest.ts
								  JOIN merchant m
								    ON l.merchant_id=m.id
									  JOIN request r
										ON r.merchant_id=m.id 
											JOIN user u
												ON u.id=r.user_id
										WHERE r.user_id=? AND active=TRUE`, userId)
	defer r.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var merchant models.MerchantInfo
	var location string
	if r.Next() {
		err = r.Scan(&merchant.Online, &merchant.Email, &merchant.Name, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Image, &merchant.MerchantID, &location, &merchant.Distance, &merchant.Accepted)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		merchant.Location.Long, merchant.Location.Lat, _ = getLongLat(location)

		var data models.RequestData
		data.UserId = userId
		data.Data = merchant.MerchantID

		merchant.IsFavorite, _ = isFavorite(data)
		return merchant, nil
	}
	return nil, nil
}

func UpdateRequest(req models.RequestEntity) (err error) {
	location:= fmt.Sprintf("POINT(%f %f)", req.Location.Long, req.Location.Lat)
	_, err = DB.Exec(`UPDATE request SET user_id=?, merchant_id=?, location=ST_GeomFromText(?), comment=?, accepted=?, active=? WHERE id=?`, req.UserID, req.MerchantID, location, req.Comment, req.Accepted, req.Active, req.ID)

	return err
}

func UpdateRequestAccepted(request models.RequestEntity, claims models.JwtClaims) (err error) {
	if request.Accepted == 1 {
		_, err := CreateNotificationByUserId(request.UserID, "",  notificationRequestAcceptedMessage, claims, request)
		if err != nil {
			fmt.Println(err)
		}

		err = InsertUserNotification(request.UserID, 1, int64(request.ID), notificationRequestAcceptedMessage)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := CreateNotificationByUserId(request.UserID, "",  notificationRequestDeclinedMessage,claims, request)
		if err != nil {
			fmt.Println(err)
		}

		err = InsertUserNotification(request.UserID, 1, int64(request.ID), notificationRequestDeclinedMessage)

		//RemoveRequestsByUserID(request.UserID)
		//return nil
	}

	_, err = DB.Exec(`UPDATE request SET accepted=? WHERE id=?`, request.Accepted, request.ID)
	return err
}

func GetRequests() (interface{}, error) {
	r, err := DB.Query(`SELECT user_id, merchant_id, ST_ASTEXT(location) FROM request`)
	defer r.Close()

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
	_, err = DB.Exec(`UPDATE request SET active=FALSE WHERE user_id=?`, userID)
	_, err = DB.Exec(`# DELETE FROM request WHERE user_id=?`, userID)
	return err
}
