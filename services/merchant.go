package services

import (
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/db"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var DB *sqlx.DB

func init() {
	DB, _ = db.Config()
}

func GetMerchants() (merchants []models.Merchant, err error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT * FROM merchant`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		r.Scan(&merchant.ID, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
		merchants = append(merchants, merchant)
	}
	return merchants, err
}

func GetMerchantByEmail(email string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT * FROM merchant WHERE email=?`, email)
	if err != nil {
		return nil, err
	}

	if r.Next() {
		r.Scan(&merchant.ID, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
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
		r.Scan(&merchant.ID, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image)
		return merchant, nil
	}
	return nil, nil
}

func CreateMerchant(merchant models.Merchant) (models.Merchant, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(merchant.Password), bcrypt.DefaultCost)
	if err != nil {
		return merchant, err
	}

	merchant.Password = string(password)
	uid, _ := uuid.NewV4()
	merchant.ID = uid.String()

	_, err = DB.Exec(`INSERT INTO merchant (id, phone_number, email, name, password) VALUES (?, ?, ?, ?, ?)`, merchant.ID, merchant.PhoneNumber, merchant.Email, merchant.Name, merchant.Password)

	if err != nil {
		return merchant, errors.New("email exists")
	}

	return merchant, nil
}

func AddNewLocation(location models.Location) (error) {
	point := fmt.Sprintf(`POINT(%f %f)`, location.Location.Lat, location.Location.Long)
	_, err:= DB.Exec(`INSERT INTO location (merchant_id, location) VALUES (?, ST_GeomFromText(?))`, location.MerchantID, point)
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
		err = r.Scan(&location.MerchantID, &point)
		if err != nil {
			return nil, err
		}
		location.Location.Lat, location.Location.Long, _ = getLatLong(point)

		return location, nil
	}
	return nil, nil
}

func GetNearMerchantsLastLocation(location models.Location) {

}
