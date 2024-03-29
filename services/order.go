package services

import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"log"
)

type (
	OrderMerchantEntity struct {
		ID         int         `json:"id"`
		User       models.User `json:"user"`
		MerchantID string      `json:"merchant_id"`
		StripeID   string      `json:"stripe_id"`
		Refund     bool        `json:"refund"`
		Amount     float64     `json:"amount"`
		Date       string      `json:"date"`
		Items      interface{} `json:"items"`
	}

	OrderUserEntity struct {
		ID       int             `json:"id"`
		UserId   string          `json:"user_id"`
		Merchant models.Merchant `json:"merchant"`
		StripeId string          `json:"stripe_id"`
		Refund   bool            `json:"refund"`
		Amount   float64         `json:"amount"`
		Date     string          `json:"date"`
		Items    interface{}     `json:"items"`
	}
)

func CreateOrder(data models.OrderRequestData) (interface{}, error) {
	r, err := DB.Exec(`INSERT INTO m_order (user_id, merchant_id, stripe_id, amount) VALUES (?,?,?,?)`, data.UserID, data.MerchantID, data.StripeID, data.Amount)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	orderId, _ := r.LastInsertId()

	for _, item := range data.Items {
		_, err = DB.Exec(`INSERT INTO m_order_product (order_id, product_id, quantity, price) VALUES (?,?,?,?)`, orderId, item.Product.ID, item.Quantity, item.Price)

		if err != nil {
			log.Println(err)
		}
	}

	return orderId, err
}

func getItemOrder(orderId int) (items []interface{}, err error) {
	r, err := DB.Query(`SELECT o.quantity, o.price, p.id, p.merchant_id, p.name, p.price, p.image 
								FROM m_order_product o
								LEFT JOIN product p on o.product_id = p.id
								WHERE order_id=?`, orderId)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	type Item struct {
		Product  models.Product `json:"product"`
		Quantity int            `json:"quantity"`
		Price    float64        `json:"price"`
	}

	for r.Next() {
		var item Item
		_ = r.Scan(&item.Quantity, &item.Price, &item.Product.ID, &item.Product.MerchantId, &item.Product.Name, &item.Product.Price, &item.Product.Image)
		items = append(items, item)
	}

	return items, nil
}

func GetOrderHistoryByUserId(userID string) (orders []interface{}, err error) {
	r, err := DB.Query(`SELECT o.id, user_id, merchant_id, o.stripe_id, refund, amount, date , online, mobile, phone_number, email, name, ST_AsText(last_location), image 
								FROM m_order o
								LEFT JOIN merchant m ON o.merchant_id=m.id 
								WHERE user_id=?
								ORDER BY date DESC 
								LIMIT 20
								`, userID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer r.Close()

	for r.Next() {
		var order OrderUserEntity
		var location string
		_ = r.Scan(&order.ID, &order.UserId, &order.Merchant.ID, &order.StripeId, &order.Refund, &order.Amount, &order.Date,
			&order.Merchant.Online, &order.Merchant.Mobile, &order.Merchant.PhoneNumber, &order.Merchant.Email, &order.Merchant.Name, &location, &order.Merchant.Image)
		order.Merchant.LastLocation.Long, order.Merchant.LastLocation.Lat, _ = getLongLat(location)
		order.Items, _ = getItemOrder(order.ID)
		orders = append(orders, order)
	}

	return orders, nil
}

func GetOrderById(id int) (order interface{}, err error) {

	r, err := DB.Query(`SELECT o.id, user_id, merchant_id, o.stripe_id, refund, amount, date , online, mobile, phone_number, email, name, ST_AsText(last_location), image 
								FROM m_order o
								LEFT JOIN merchant m ON o.merchant_id=m.id 
								WHERE o.id=?`, id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer r.Close()

	if r.Next() {
		var order OrderUserEntity
		var location string
		_ = r.Scan(&order.ID, &order.UserId, &order.Merchant.ID, &order.StripeId, &order.Refund, &order.Amount, &order.Date,
			&order.Merchant.Online, &order.Merchant.Mobile, &order.Merchant.PhoneNumber, &order.Merchant.Email, &order.Merchant.Name, &location, &order.Merchant.Image)
		order.Merchant.LastLocation.Long, order.Merchant.LastLocation.Lat, _ = getLongLat(location)
		order.Items, _ = getItemOrder(order.ID)

		return order, nil
	}

	return nil, nil
}


func GetOrderPaymentByMerchantId(merchantID string) (orders []interface{}, err error) {
	r, err := DB.Query(`SELECT o.id, user_id, merchant_id, stripe_id, refund, amount, date, phone_number, email, name, image, ST_AsText(last_location)
								FROM m_order o
								LEFT JOIN user u ON o.user_id=u.id
								WHERE merchant_id=?`, merchantID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer r.Close()

	for r.Next() {
		var order OrderMerchantEntity
		var location string
		err = r.Scan(&order.ID, &order.User.ID, &order.MerchantID, &order.StripeID, &order.Refund, &order.Amount, &order.Date,
			&order.User.PhoneNumber, &order.User.Email, &order.User.Name, &order.User.Image, &location)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		order.User.LastLocation.Long, order.User.LastLocation.Lat, _ = getLongLat(location)
		order.Items, _ = getItemOrder(order.ID)

		orders = append(orders, order)
	}

	return orders, err
}

func GetOrderPaymentById(orderId int) (order interface{}, err error) {
	r, err := DB.Query(`SELECT o.id, user_id, merchant_id, stripe_id, refund, amount, date, phone_number, email, name, image, ST_AsText(last_location)
								FROM m_order o
								LEFT JOIN user u ON o.user_id=u.id
								WHERE o.id=?`, orderId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer r.Close()

	if r.Next() {
		var order OrderMerchantEntity
		var location string
		err = r.Scan(&order.ID, &order.User.ID, &order.MerchantID, &order.StripeID, &order.Refund, &order.Amount, &order.Date,
			&order.User.PhoneNumber, &order.User.Email, &order.User.Name, &order.User.Image, &location)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		order.User.LastLocation.Long, order.User.LastLocation.Lat, _ = getLongLat(location)
		order.Items, _ = getItemOrder(order.ID)

		return order, nil
	}

	return nil, err
}

func GetOrderEntityById(orderId int) (order interface{}, err error) {
	r, err := DB.Query(`SELECT id, user_id, merchant_id, stripe_id, refund, amount, date FROM m_order WHERE id=?`, orderId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer r.Close()

	if r.Next() {
		var order models.OrderEntity
		err = r.Scan(&order.ID, &order.UserId, &order.MerchantId, &order.StripeId, &order.Refund, &order.Amount, &order.Date)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return order, nil
	}
	return nil, nil
}

func UpdateOrder(order models.OrderEntity) error {
	_, err := DB.Exec(`UPDATE m_order SET refund=? WHERE id=?`, order.Refund, order.ID)
	return err
}
