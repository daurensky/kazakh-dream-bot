package db

import (
	"github.com/daurensky/kazakh-dream-bot/models"
	"github.com/lib/pq"
)

func GetClientProducts(telegramId int64, withDuplicates bool) ([]models.Product, error) {
	var products []models.Product
	var sql string

	db, err := Connect()

	if err != nil {
		return nil, err
	}

	if withDuplicates {
		sql = `
			SELECT p.id,
				   p.price,
				   p.photo_url,
				   p.composition,
			       p.name
			FROM kazakh_dream.public.cart c
				INNER JOIN kazakh_dream.public.products p on p.id = c.product_id
			WHERE c.telegram_id = $1
		`
	} else {
		sql = `
			SELECT p.id,
			   	   SUM(p.price),
				   p.photo_url,
				   p.composition,
				   CASE WHEN COUNT(p.id) > 1 THEN CONCAT(p.name, ' ', COUNT(p.id), ' шт.') ELSE p.name END AS name
			FROM kazakh_dream.public.cart c
				INNER JOIN kazakh_dream.public.products p on p.id = c.product_id
			WHERE c.telegram_id = $1
			GROUP BY p.id
		`
	}
	productsFromDB, err := db.Query(sql, telegramId)

	for productsFromDB.Next() {
		product := models.Product{}

		err := productsFromDB.Scan(&product.Id, &product.Price, &product.PhotoUrl, pq.Array(&product.Composition), &product.Name)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
