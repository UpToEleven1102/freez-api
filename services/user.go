package services

import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)


func CreateUser(user models.User) (interface{}, error) {
	uid, _ := uuid.NewV4()
	user.ID = uid.String()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	location := fmt.Sprintf("POINT(%f %f)", user.LastLocation.Long, user.LastLocation.Lat)
	_, err := DB.Exec(`INSERT INTO user (id, phone_number, email, name, password, last_location, image) VALUES(?,?,?,?,?,ST_GeomFromText(?),?)`, user.ID, user.PhoneNumber, user.Email, user.Name, user.Password, location, user.Image);
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(user models.User) (err error) {
	_, err = DB.Exec(`UPDATE user SET phone_number=?,email=?,name=?,image=? WHERE id=?;`,user.PhoneNumber,user.Email,user.Name,user.Image,user.ID)
	_, err = DB.Exec(`UPDATE m_option SET notif_fav_nearby=? WHERE user_id=?;`, user.Option.NotifFavNearby, user.ID)
	return err
}

func GetUserByEmail(email string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, phone_number, email, name, password, image, ST_AsText(last_location) FROM user WHERE email=?`, email)
	defer r.Close()

	if err != nil {
		return nil, err
	}
	var location string
	var user models.User
	if r.Next() {
		r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location)
		return user, nil
	}

	user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)

	return nil, nil
}

func GetUserById(id string) (interface{}, error) {
	r, err := DB.Query(`SELECT u.id, phone_number, email, name, password, image, ST_AsText(last_location), notif_fav_nearby 
								  FROM user u 
								    INNER JOIN m_option o 
								      ON u.id=o.user_id WHERE u.id=?`, id)
	defer r.Close()

	if err != nil {
		return nil, err
	}

	var location string
	var user models.User
	if r.Next() {
		r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location, &user.Option.NotifFavNearby)
		user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)
		return user, nil
	}

	return nil, nil
}

func GetUserByPhoneNumber(phoneNumber string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, phone_number, email, name, password, image, ST_AsText(last_location) FROM user WHERE phone_number=?`, phoneNumber)
	defer r.Close()

	if err != nil {
		return nil, err
	}
	var location string
	var user models.User
	if r.Next() {
		r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location)
		return user, nil
	}

	user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)

	return nil, nil
}

func UpdateUserLocation(user models.User) (interface{}, error) {
	point:= fmt.Sprintf(`POINT(%f %f)`,user.LastLocation.Long, user.LastLocation.Lat)

	_, err := DB.Exec(`UPDATE user SET last_location=ST_GeomFromText(?) WHERE id=?`, point, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AddFavorite(data models.RequestData) (err error) {
	_, err = DB.Exec(`INSERT INTO favorite (user_id, merchant_id) VALUES(?,?)`, data.UserId, data.Data)
	return err
}

func RemoveFavorite(data models.RequestData) (err error) {
	_, err = DB.Exec(`DELETE FROM favorite WHERE user_id=? AND merchant_id=?`, data.UserId, data.Data)
	return err
}

func isFavorite(data models.RequestData) (bool, error) {
	r, err := DB.Query(`SELECT * FROM favorite WHERE user_id=? AND merchant_id=?`, data.UserId, data.Data)
	defer r.Close()

	if err != nil {
		return false, err
	}

	if r.Next() {
		return true, nil
	}
	return false, nil
}

func GetFavorites(userID string) (merchants []interface{}, err error) {
	r, err := DB.Query(`SELECT online, merchant_id, ST_AsText(m.last_location) as location, ST_Distance_Sphere(u.last_location, m.last_location) as distance, m.name, m.phone_number, m.email, mobile, m.image
								FROM favorite f
								  INNER JOIN user u
								  	ON f.user_id=u.id
									INNER JOIN merchant m
									  ON f.merchant_id=m.id
										WHERE f.user_id=?`, userID)
	defer r.Close()
	if err != nil {
		return nil, err
	}

	var merchant models.MerchantInfo
	var location string
	for r.Next() {
		err = r.Scan(&merchant.Online, &merchant.MerchantID, &location, &merchant.Distance, &merchant.Name, &merchant.PhoneNumber, &merchant.Email, &merchant.Mobile, &merchant.Image)
		if err != nil {
			return nil, err
		}
		merchant.Location.Long, merchant.Location.Lat, _ = getLongLat(location)
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}
