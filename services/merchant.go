package services

import (
	"errors"
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"github.com/satori/go.uuid"
	"github.com/stripe/stripe-go"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//func GetMerchants() (merchants []models.Merchant, err error) {
//	var merchant models.Merchant
//	var location string
//
//	r, err := DB.Query(`SELECT * FROM merchant`)
//	defer r.Close()
//
//	if err != nil {
//		return nil, err
//	}
//
//	for r.Next() {
//		r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email, &merchant.Name, &merchant.Password, &merchant.Image, &location)
//
//		merchant.LastLocation.Long, merchant.LastLocation.Lat, _ = getLongLat(location)
//
//		merchants = append(merchants, merchant)
//	}
//
//	return merchants, err
//}

func GetUserIDNotifyMerchantNearbyByMerchantID(merchantLocation models.Location) (ids []interface{}, err error) {
	location := fmt.Sprintf("POINT(%f %f)", merchantLocation.Location.Long, merchantLocation.Location.Lat)
	r, err := DB.Query(`SELECT fav.user_id 
								FROM favorite fav 
								  LEFT JOIN m_option o 
								    ON fav.user_id=o.user_id
										LEFT JOIN user u on o.user_id = u.id
								WHERE o.notif_fav_nearby=TRUE AND ST_Distance_Sphere(u.last_location, ST_GeomFromText(?)) < ? AND fav.merchant_id=?`, location, minNotifyDistance, merchantLocation.Id)

	if err != nil {
		return nil, err
	}

	for r.Next() {
		var userId string
		err = r.Scan(&userId)

		if err != nil {
			return nil, err
		}
		ids = append(ids, userId)
	}
	return ids, nil
}

func CreateMerchantFB(merchant models.Merchant) (models.Merchant, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(merchant.Password), bcrypt.DefaultCost)
	if err != nil {
		return merchant, err
	}

	merchant.Password = string(password)
	uid, _ := uuid.NewV4()
	merchant.ID = uid.String()

	location := fmt.Sprintf("POINT(%f %f)", merchant.LastLocation.Long, merchant.LastLocation.Lat)

	_, err = DB.Exec(`INSERT INTO merchant (id, mobile, phone_number, email, name, description, password, last_location, stripe_id, image, facebook_id, category) VALUES (?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?, ?, ?, ?)`, merchant.ID, merchant.Mobile, merchant.PhoneNumber, merchant.Email, merchant.Name, merchant.Description, merchant.Password, location, merchant.StripeID, merchant.Image, merchant.FacebookID, merchant.Category)

	if err != nil {
		return merchant, err
	}

	return merchant, nil
}

func CreateMerchant(merchant models.Merchant) (models.Merchant, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(merchant.Password), bcrypt.DefaultCost)
	if err != nil {
		return merchant, err
	}

	merchant.Password = string(password)
	uid, _ := uuid.NewV4()
	merchant.ID = uid.String()

	location := fmt.Sprintf("POINT(%f %f)", merchant.LastLocation.Long, merchant.LastLocation.Lat)

	_, err = DB.Exec(`INSERT INTO merchant (id, mobile, phone_number, email, name, description, password, last_location, stripe_id, category) VALUES (?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?, ?)`, merchant.ID, merchant.Mobile, merchant.PhoneNumber, merchant.Email, merchant.Name, merchant.Description, merchant.Password, location, merchant.StripeID, merchant.Category)

	if err != nil {
		return merchant, err
	}

	return merchant, nil
}

func GetMerchantStripeIdByMerchantId(merchantId string) (id string, err error) {
	r, err := DB.Query(`SELECT stripe_id FROM merchant WHERE id=?`, merchantId)

	if err != nil {
		panic(err)
	}

	if r.Next() {
		err = r.Scan(&id)
		if err != nil {
			panic(err)
		}
	}

	return id, err
}

func FilterMerchantByName(data models.SearchData, location models.Location) (merchants []interface{}, err error) {

	point := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)

	r, err := DB.Query(`SELECT online, category, email, name, description, mobile, phone_number, image, l.merchant_id, ST_AsText(location) as location, ST_Distance_Sphere(location, ST_GeomFromText(?)) as distance
								FROM location l INNER JOIN (
								    SELECT merchant_id, MAX(ts) AS ts FROM location GROUP BY merchant_id
								  ) latest
								  ON l.ts=latest.ts
								  LEFT JOIN merchant m
								    ON l.merchant_id=m.id
									  WHERE INSTR(name, ?) > 0
									  LIMIT ?
									  `, point, data.SearchText, data.Limit)

	if err != nil {
		return nil, err
	}
	defer r.Close()
	var loc string

	for r.Next() {
		var merchant models.MerchantInfo
		_ = r.Scan(&merchant.Online, &merchant.Category, &merchant.Email, &merchant.Name, &merchant.Description, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Image, &merchant.MerchantID, &loc, &merchant.Distance)

		merchant.Location.Long, merchant.Location.Lat, err = getLongLat(loc)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		merchants = append(merchants, merchant)
	}

	return merchants, err
}

func GetMerchantStripeAccount(merchantId string) (*stripe.Account, error) {
	stripeId, err := GetMerchantStripeIdByMerchantId(merchantId)

	if err != nil {
		panic(err)
	}

	return StripeConnectGetAccountById(stripeId)
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

	merchant.Online = !merchant.Online

	_, err = DB.Exec(`UPDATE merchant SET online=? WHERE id=?;`, merchant.Online, merchant.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetMerchantByEmail(email string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT id, online, mobile, phone_number, email, name, description, password, image, ST_AsText(last_location) FROM merchant WHERE email=?`, email)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string

	if r.Next() {
		err = r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email,
			&merchant.Name, &merchant.Description, &merchant.Password, &merchant.Image, &location)
		merchant.LastLocation.Long, merchant.LastLocation.Lat, _ = getLongLat(location)
		if err != nil {
			return nil, err
		}
		return merchant, nil
	}
	return nil, nil
}

func GetMerchantById(id string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT id, online, mobile, category, phone_number, email, name, description, password, image, ST_AsText(last_location) FROM merchant WHERE id=?`, id)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string

	if r.Next() {
		err = r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.Category, &merchant.PhoneNumber,
			&merchant.Email, &merchant.Name, &merchant.Description, &merchant.Password, &merchant.Image, &location)
		merchant.LastLocation.Long, merchant.LastLocation.Lat, _ = getLongLat(location)
		if err != nil {
			return nil, err
		}
		return merchant, nil
	}
	return nil, nil
}

func GetMerchantByPhoneNumber(phoneNumber string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT id, online, mobile, phone_number, email, name, description, password, image, ST_AsText(last_location) FROM merchant WHERE phone_number=?`, phoneNumber)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string

	if r.Next() {
		err = r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email,
			&merchant.Name, &merchant.Description, &merchant.Password, &merchant.Image, &location)
		if err != nil {
			return nil, err
		}
		merchant.LastLocation.Long, merchant.LastLocation.Lat, _ = getLongLat(location)
		return merchant, nil
	}
	return nil, nil
}

func GetMerchantByFacebookID(facebookID string) (interface{}, error) {
	var merchant models.Merchant

	r, err := DB.Query(`SELECT id, online, mobile, phone_number, email, name, description, password, image, ST_AsText(last_location) FROM merchant WHERE facebook_id=?`, facebookID)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string

	if r.Next() {
		err = r.Scan(&merchant.ID, &merchant.Online, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Email,
			&merchant.Name, &merchant.Description, &merchant.Password, &merchant.Image, &location)
		if err != nil {
			return nil, err
		}
		merchant.LastLocation.Long, merchant.LastLocation.Lat, _ = getLongLat(location)
		return merchant, nil
	}
	return nil, nil
}

func UpdateMerchant(merchant models.Merchant) (err error) {
	_, err = DB.Exec(`UPDATE merchant SET mobile=?,phone_number=?,email=?,name=?,description=?, image=? WHERE id=?;`, merchant.Mobile, merchant.PhoneNumber, merchant.Email, merchant.Name, merchant.Description, merchant.Image, merchant.ID)
	return err
}

func UpdateFoodType(merchant models.Merchant) (err error) {
	_, err = DB.Exec("UPDATE merchant SET category=? WHERE id=?", merchant.Category, merchant.ID)
	return err
}

func AddNewLocation(location models.Location) (error) {
	point := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)
	_, err := DB.Exec(`INSERT INTO location (merchant_id, location) VALUES (?, ST_GeomFromText(?))`, location.Id, point)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`UPDATE merchant SET last_location=ST_GeomFromText(?) WHERE id=?`, point, location.Id)
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
	defer r.Close()

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

func GetMerchantInfoById(id string, location models.Location) (interface{}, error) {
	userLocation := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)
	r, err := DB.Query(`SELECT online, email, name, description, mobile, phone_number, image, category, l.merchant_id, ST_AsText(location) as location, ST_Distance_Sphere(location, ST_GeomFromText(?)) as distance
								FROM location l INNER JOIN (
								    SELECT merchant_id, MAX(ts) AS ts FROM location GROUP BY merchant_id
								  ) latest
								  ON l.ts=latest.ts
								  JOIN merchant m
								    ON l.merchant_id=m.id
									  WHERE l.merchant_id=?`, userLocation, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer r.Close()

	var loc string

	if r.Next() {
		var merchant models.MerchantInfo
		_ = r.Scan(&merchant.Online, &merchant.Email, &merchant.Name, &merchant.Description, &merchant.Mobile, &merchant.PhoneNumber, &merchant.Image, &merchant.Category, &merchant.MerchantID, &loc, &merchant.Distance)

		merchant.Location.Long, merchant.Location.Lat, err = getLongLat(loc)
		if err != nil {
			return nil, err
		}

		var data models.RequestData
		data.UserId = location.Id
		data.Data = merchant.MerchantID

		merchant.IsFavorite, _ = isFavorite(data)

		merchant.Products, _ = GetProducts(merchant.MerchantID)
		return merchant, nil
	}

	return nil, nil
}

func GetNearbyMerchantsLastLocation(location models.Location, filters ...string) (merchants []interface{}, err error) {
	userLocation := fmt.Sprintf(`POINT(%f %f)`, location.Location.Long, location.Location.Lat)

	var filter string

	if len(filters) > 0 {
		filter = ` AND (`
		for idx, f := range filters {
			filter += `category='` + f + `'`
			if idx != len(filters)-1 {
				filter += ` OR `
			}
		}
		filter += ")"
	}

	r, err := DB.Query(`SELECT online, category, email, name, description, mobile, phone_number, image, l.merchant_id, ST_AsText(location) as location, ST_Distance_Sphere(location, ST_GeomFromText(?)) as distance
								FROM location l INNER JOIN (
								    SELECT merchant_id, MAX(ts) AS ts FROM location GROUP BY merchant_id
								  ) latest
								  ON l.ts=latest.ts
								  LEFT JOIN merchant m
								    ON l.merchant_id=m.id
									  HAVING distance < ? AND online=true`+filter, userLocation, minDistance)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var loc string

	for r.Next() {
		var merchant models.MerchantInfo
		_ = r.Scan(&merchant.Online, &merchant.Category, &merchant.Email, &merchant.Name, &merchant.Description,
			&merchant.Mobile, &merchant.PhoneNumber, &merchant.Image, &merchant.MerchantID, &loc, &merchant.Distance)

		merchant.Location.Long, merchant.Location.Lat, err = getLongLat(loc)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		var data models.RequestData
		data.UserId = location.Id
		data.Data = merchant.MerchantID

		merchant.IsFavorite, _ = isFavorite(data)
		merchant.Products, _ = GetProducts(merchant.MerchantID)

		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

//stripe operations
func GetStripeCardList(merchantId string) ([]*stripe.Card, error) {
	id, err := GetMerchantStripeIdByMerchantId(merchantId)

	if err != nil {
		return nil, err
	}

	return StripeConnectGetCardListByStripeId(id)
}
