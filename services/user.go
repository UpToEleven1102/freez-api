package services


import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CreateUser(user models.User) (interface{}, error) {
	uid, _ := uuid.NewV4()
	user.ID = uid.String()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	location := fmt.Sprintf("POINT(%f %f)", user.LastLocation.Long, user.LastLocation.Lat)
	_, err := DB.Exec(`INSERT INTO user (id, phone_number, email, name, password, last_location) VALUES(?,?,?,?,?,ST_GeomFromText(?))`, user.ID, user.PhoneNumber, user.Email, user.Name, user.Password, location);
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUserFB(user models.User) (interface{}, error) {
	uid, _ := uuid.NewV4()
	user.ID = uid.String()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	location := fmt.Sprintf("POINT(%f %f)", user.LastLocation.Long, user.LastLocation.Lat)
	_, err := DB.Exec(`INSERT INTO user (id, email, phone_number, name, password, last_location, image, facebook_id) VALUES(?,?,?,?,?,ST_GeomFromText(?),?,?)`, user.ID, user.Email, user.PhoneNumber, user.Name, user.Password, location, user.Image, user.FacebookID);
	if err != nil {
		return nil, err
	}

	return user, nil
}


func UpdateUser(user models.User) (err error) {
	_, err = DB.Exec(`UPDATE user SET phone_number=?,email=?,name=?,image=? WHERE id=?;`, user.PhoneNumber, user.Email, user.Name, user.Image, user.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`UPDATE m_option SET notif_fav_nearby=? WHERE user_id=?;`, user.Option.NotifFavNearby, user.ID)
	return err
}

func GetUserByEmail(email string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, phone_number, email, name, password, image, ST_AsText(last_location), freez_point
								FROM user WHERE email=?`, email)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string
	var user models.User
	if r.Next() {
		err = r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location, &user.FreezPoint)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)

	return nil, nil
}

func GetUserById(id string) (interface{}, error) {
	r, err := DB.Query(`SELECT u.id, phone_number, email, name, password, image, ST_AsText(last_location), notif_fav_nearby, freez_point
								  FROM user u 
								    INNER JOIN m_option o 
								      ON u.id=o.user_id WHERE u.id=?`, id)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	var location string
	var user models.User
	if r.Next() {
		err = r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location, &user.Option.NotifFavNearby, &user.FreezPoint)
		if err != nil {
			return nil, err
		}
		user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)
		return user, nil
	}

	return nil, nil
}

func GetUserByPhoneNumber(phoneNumber string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, phone_number, email, name, password, image, ST_AsText(last_location), freez_point 
								FROM user WHERE phone_number=?`, phoneNumber)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string
	var user models.User
	if r.Next() {
		err = r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location, &user.FreezPoint)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)

	return nil, nil
}

func UpdateUserLocation(user models.User) (interface{}, error) {
	point := fmt.Sprintf(`POINT(%f %f)`, user.LastLocation.Long, user.LastLocation.Lat)

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
	defer r.Close()

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
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var merchant models.MerchantInfo
	var location string
	for r.Next() {
		err = r.Scan(&merchant.Online, &merchant.MerchantID, &location, &merchant.Distance, &merchant.Name, &merchant.PhoneNumber, &merchant.Email, &merchant.Mobile, &merchant.Image)
		if err != nil {
			return nil, err
		}
		merchant.Location.Long, merchant.Location.Lat, _ = getLongLat(location)
		merchant.IsFavorite = true
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func ChargeUser(data models.OrderRequestData) (orderId interface{}, err error) {
	stripeAccId, err := GetMerchantStripeIdByMerchantId(data.MerchantID)

	if err != nil {
		return nil, err
	}

	res, err := StripeConnectDestinationCharge(data.StripeToken, stripeAccId, "Testing", data.Amount)

	log.Println(res)
	if err != nil {
		return nil, err
	}

	data.StripeID = res.ID

	return CreateOrder(data)
}

func AddPointsPerPurchase(userID string, numPoints int) error {
	_, err := DB.Exec(`UPDATE user SET freez_point=freez_point+? WHERE id=?`, numPoints, userID)

	return err
}

func GetUserByFbId(facebookID string) (interface{}, error) {
	r, err := DB.Query(`SELECT id, phone_number, email, name, password, image, ST_AsText(last_location), freez_point
			FROM user WHERE facebook_id=?`, facebookID)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var location string
	var user models.User
	if r.Next() {
		err = r.Scan(&user.ID, &user.PhoneNumber, &user.Email, &user.Name, &user.Password, &user.Image, &location, &user.FreezPoint)

		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user.LastLocation.Long, user.LastLocation.Lat, _ = getLongLat(location)

	return nil, nil
}