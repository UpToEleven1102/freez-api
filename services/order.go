package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"log"
)

func CreateOrder(data models.OrderRequestData) error {
	r, err := DB.Exec(`INSERT INTO m_order (user_id, merchant_id, stripe_id) VALUES (?,?,?)`, data.UserID, data.MerchantID, data.StripeID)

	if err != nil {
		log.Println(err)
		return err
	}

	orderId, _ := r.LastInsertId()

	for _, item := range data.Items {
		_, err = DB.Exec(`INSERT INTO m_order_product (order_id, product_id, quantity, price) VALUES (?,?,?,?)`, orderId, item.Product.ID, item.Quantity, item.Price)

		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func getItemOrder(orderId int) (items []interface{}, err error) {
	r, err := DB.Query(`SELECT o.quantity, o.price, p.id, p.merchant_id, p.name, p.price, p.image 
								FROM m_order_product o
								LEFT JOIN product p on o.product_id = p.id
								WHERE order_id=?`, orderId)
	defer r.Close()
	if err != nil {
		log.Println(err)
		return items, err
	}

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
	r, err := DB.Query(`SELECT id, user_id, merchant_id, stripe_id, refund, amount FROM m_order WHERE user_id=?`, userID)

	defer r.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for r.Next() {
		var order models.OrderEntity
		_ = r.Scan(&order.ID, &order.UserId, &order.MerchantId, &order.StripeId, &order.Refund, &order.Amount)
		order.Items, _ = getItemOrder(order.ID)
		orders = append(orders, order)
	}

	return orders, nil
}
