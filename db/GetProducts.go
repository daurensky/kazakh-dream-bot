package db

import (
	"github.com/daurensky/kazakh-dream-bot/models"
	"github.com/lib/pq"
)

func GetProducts() ([]models.Product, error) {
	var products []models.Product

	db, err := Connect()

	if err != nil {
		return nil, err
	}

	productsFromDB, err := db.Query("SELECT * FROM kazakh_dream.public.products")

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
