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

	_, err := DB.Exec(`INSERT INTO user (id, phone_number, email, name, password, image) VALUES(?,?,?,?,?,?)`, user.ID, user.PhoneNumber, user.Email, user.Name, user.Password, user.Image);
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (interface{}, error) {
	r, err := DB.Query(`SELECT * FROM user WHERE email=?`, email)

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
	r, err := DB.Query(`SELECT * from user WHERE id=?`, id)

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
	if err != nil {
		return false, err
	}

	if r.Next() {
		return true, nil
	}
	return false, nil
}
//
//func GetFavorites(user_id string) (err error) {
//	_, err = DB.Query(`SELECT online, email, name, mobile, phone_number, image, merchant_id
//								FROM favorite f
//								  JOIN merchant m
//								    ON f.merchant_id=m.id GROUP BY merchant_id`)
//}