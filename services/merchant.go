package services

import (
	"golang.org/x/crypto/bcrypt"
	"../models"
	"github.com/jmoiron/sqlx"
	"../db"
	"github.com/satori/go.uuid"
	"errors"
)

var DB *sqlx.DB

type _response struct {
	data interface{}
}

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

func GetMerchantById(id int64) (interface{}, error) {
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
