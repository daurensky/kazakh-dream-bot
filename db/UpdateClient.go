package db

import "github.com/daurensky/kazakh-dream-bot/models"

func UpdateClient(client models.Client) error {
	db, err := Connect()

	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE kazakh_dream.public.clients SET name = $1, phone = $2, address = $3 WHERE telegram_id = $4",
		client.Name,
		client.Phone,
		client.Address,
		client.TelegramId,
	)

	if err != nil {
		return err
	}

	return nil
}
