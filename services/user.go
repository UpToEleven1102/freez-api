package services

import (
	"../models"
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
)
func CreateUser(user models.User) (interface{}, error) {
	uid, _ := uuid.NewV4()
	user.ID = uid.String()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	_, err := DB.Exec(`INSERT INTO user (id, phone_number, email, name, password) VALUES(?,?,?,?,?)`, user.ID, user.PhoneNumber, user.Email, user.Name, user.Password);
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

	var user models.User
	if r.Next() {
		r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image)
		return user, nil
	}

	return nil, nil
}

func GetUserById(id string) (interface{}, error) {
	r, err := DB.Query(`SELECT * from user WHERE id=?`, id)

	if err != nil {
		return nil, err
	}

	var user models.User
	if r.Next() {
		r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image)
		return user, nil
	}

	return nil, nil
}
