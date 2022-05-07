package db

import "github.com/daurensky/kazakh-dream-bot/models"

func StoreOrder(client models.Client) (int, error) {
	var orderId int

	db, err := Connect()

	if err != nil {
		return orderId, err
	}

	err = db.QueryRow(
		"INSERT INTO kazakh_dream.public.orders (client_id) VALUES ($1) RETURNING id",
		client.TelegramId,
	).Scan(&orderId)

	if err != nil {
		return orderId, err
	}

	return orderId, nil
}
