package db

import (
	"github.com/daurensky/kazakh-dream-bot/models"
)

func StoreClient(client models.Client) error {
	db, err := Connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO kazakh_dream.public.clients (telegram_id, name, phone, address) VALUES ($1, $2, $3, $4)",
		client.TelegramId,
		client.Name,
		client.Phone,
		client.Address,
	)

	if err != nil {
		return err
	}

	return nil
}
