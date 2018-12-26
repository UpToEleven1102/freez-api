package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"log"
)

func CreateProduct(product models.Product) error {
	_, err := DB.Exec(`INSERT INTO product (name, price, merchant_id) VALUES(?, ?, ?)`, product.Name, product.Price, product.MerchantId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetProducts(merchantID string) (products []interface{}, err error) {
	r, err := DB.Query(`SELECT id, name, price, merchant_id, image FROM product WHERE merchant_id=?`, merchantID)
	defer r.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var product models.Product
	for r.Next() {
		err = r.Scan(&product.ID, &product.Name, &product.Price, &product.MerchantId, &product.Image)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func UpdateProduct(product models.Product) error {
	_, err := DB.Exec(`UPDATE product SET name=?, price=? WHERE id=?`, product.Name, product.Price, product.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteProduct(product models.Product) error {
	_, err := DB.Exec(`DELETE FROM product WHERE id=?`, product.ID)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}