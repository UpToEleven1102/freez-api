package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetMerchants() (merchants []models.Merchant, err error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT * FROM merchant`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
		merchants = append(merchants, merchant)
	}
	return merchants, err
}

func ChangeOnlineStatus(merchantId string) error {
	m, err := GetMerchantById(merchantId)
	if err != nil {
		return err
	}

	if m == nil {
		return errors.New("Unauthorized")
	}

	merchant := m.(models.Merchant)

	merchant.Online	= !merchant.Online

	_, err = DB.Exec(`UPDATE merchant SET online=? WHERE id=?;`, merchant.Online, merchant.ID)

	if err != nil {
		return err
	}

	return nil
}
func GetMerchantByEmail(email string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT * FROM merchant WHERE email=?`, email)
	if err != nil {
		return nil, err
	}

	if r.Next() {
		r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
		return merchant, nil
	}
	return nil, nil
}

func GetMerchantById(id string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT * FROM merchant WHERE id=?`, id)
	if err != nil {
		return nil, err
	}

	if r.Next() {
		r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
		return merchant, nil
	}
	return nil, nil
}

func CreateMerchant(merchant models.Merchant) (models.Merchant, error) {
	b, _ := json.Marshal(merchant)

	fmt.Println(string(b))

	password, err := bcrypt.GenerateFromPassword([]byte(merchant.Password), bcrypt.DefaultCost)
	if err != nil {
		return merchant, err
	}

	merchant.Password = string(password)
	uid, _ := uuid.NewV4()
	merchant.ID = uid.String()

	_, err = DB.Exec(`INSERT INTO merchant (id, mobile, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?, ?)`, merchant.ID, merchant.Mobile, merchant.PhoneNumber, merchant.Email, merchant.Name, merchant.Password)

	if err != nil {
		return merchant, errors.New("email exists")
	}

	return merchant, nil
}

func AddNewLocation(location models.Location) (error) {
	point := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)
	_, err:= DB.Exec(`INSERT INTO location (merchant_id, location) VALUES (?, ST_GeomFromText(?))`, location.Id, point)
	if err != nil {
		return err
	}
	return nil
}

func GetLastPositionByMerchantID(merchantID string) (interface{}, error) {
	r, err := DB.Query(`SELECT merchant_id, ST_AsText(location) FROM location WHERE merchant_id=? ORDER BY ts DESC LIMIT 1;`, merchantID)
	if err != nil {
		return nil, err
	}

	var location models.Location
	var point string
	if r.Next() {
		err = r.Scan(&location.Id, &point)
		if err != nil {
			return nil, err
		}
		location.Location.Long, location.Location.Lat, _ = getLongLat(point)

		return location, nil
	}
	return nil, nil
}

func GetNearMerchantsLastLocation(location models.Location) (merchants []interface{}, err error){
	userLocation := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)

	r, err := DB.Query(`SELECT online, email, name, mobile, phone_number, image, l.merchant_id, ST_AsText(location) as location, ST_Distance_Sphere(location, ST_GeomFromText(?))*.000621371192 as distance
								FROM location l INNER JOIN (
								    SELECT merchant_id, MAX(ts) AS ts FROM location GROUP BY merchant_id
								  ) latest
								  ON l.ts=latest.ts
								  JOIN merchant m
								    ON l.merchant_id=m.id
									  HAVING distance < 3 AND online=true`, userLocation)

	if err != nil {
		return nil, err
	}

	var loc string

	for r.Next()  {
		var merchant models.MerchantInfo
		_ = r.Scan(&merchant.Online ,&merchant.Email, &merchant.Name, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Image, &merchant.MerchantID, &loc, &merchant.Distance)

		merchant.Location.Long, merchant.Location.Lat, err = getLongLat(loc)
		if err != nil {
			return nil, err
		}

		var data models.RequestData
		data.UserId = location.Id
		data.Data = merchant.MerchantID

		merchant.IsFavorite, _ = isFavorite(data)

		merchants = append(merchants, merchant)
	}

	return merchants, nil
}
