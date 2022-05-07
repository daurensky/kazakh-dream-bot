package db

import "github.com/daurensky/kazakh-dream-bot/models"

func StoreOrderProducts(orderId int, products []models.Product) error {
	db, err := Connect()

	if err != nil {
		return err
	}

	for _, product := range products {
		_, err := db.Exec(
			"INSERT INTO kazakh_dream.public.order_product (order_id, product_id) VALUES ($1, $2)",
			orderId,
			product.Id,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
