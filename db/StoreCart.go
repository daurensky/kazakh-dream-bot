package db

import "github.com/daurensky/kazakh-dream-bot/models"

func StoreCart(telegramId int64, product models.Product) error {
	db, err := Connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO kazakh_dream.public.cart (product_id, telegram_id) VALUES ($1, $2)",
		product.Id,
		telegramId,
	)

	if err != nil {
		return err
	}

	return nil
}
