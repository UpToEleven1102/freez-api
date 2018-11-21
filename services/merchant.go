package services

import (
	"../models"
	"github.com/jmoiron/sqlx"
	"../db"
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

	for r.Next()  {
		r.Scan(&merchant.ID, &merchant.Name, &merchant.Email, &merchant.Password, &merchant.PhoneNumber, &merchant.Image)
		merchants = append(merchants, merchant)
	}
	return merchants, err
}
