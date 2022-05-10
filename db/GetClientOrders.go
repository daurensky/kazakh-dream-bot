package db

import (
	"github.com/daurensky/kazakh-dream-bot/models"
	"github.com/lib/pq"
)

func GetClientOrders(telegramId int64) ([]models.Order, error) {
	var orders []models.Order

	db, err := Connect()

	if err != nil {
		return nil, err
	}

	ordersFromDB, err := db.Query(`
		SELECT *
		FROM kazakh_dream.public.orders
		WHERE client_id = $1
		ORDER BY created_at DESC LIMIT 5
	`, telegramId)

	for ordersFromDB.Next() {
		var orderProducts []models.Product

		order := models.Order{}

		err := ordersFromDB.Scan(&order.Id, &order.Status, &order.CreatedAt, &order.ClientId)

		if err != nil {
			return nil, err
		}

		orderProductsFromDB, err := db.Query(`
			SELECT p.id,
				   SUM(p.price),
				   p.photo_url,
				   p.composition,
				   CASE WHEN COUNT(p.id) > 1 THEN CONCAT(p.name, ' ', COUNT(p.id), ' шт') ELSE p.name END AS name
			FROM kazakh_dream.public.order_product op
				INNER JOIN kazakh_dream.public.products p on p.id = op.product_id
				WHERE order_id = $1
			GROUP BY p.id
		`, order.Id)

		if err != nil {
			panic(err)
		}

		for orderProductsFromDB.Next() {
			product := models.Product{}

			err := orderProductsFromDB.Scan(
				&product.Id,
				&product.Price,
				&product.PhotoUrl,
				pq.Array(&product.Composition),
				&product.Name,
			)

			if err != nil {
				panic(err)
			}

			orderProducts = append(orderProducts, product)
		}

		order.Products = orderProducts
		orders = append(orders, order)
	}

	return orders, nil
}
