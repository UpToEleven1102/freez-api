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
