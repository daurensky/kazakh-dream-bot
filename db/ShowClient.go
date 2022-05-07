package db

import (
	"github.com/daurensky/kazakh-dream-bot/models"
)

func ShowClient(clientId int64) (models.Client, error) {
	client := models.Client{}

	db, err := Connect()

	if err != nil {
		return client, err
	}

	row := db.QueryRow("SELECT * FROM kazakh_dream.public.clients WHERE telegram_id = $1", clientId)

	err = row.Scan(&client.TelegramId, &client.Name, &client.Phone, &client.Address)

	if err != nil {
		return client, err
	}

	return client, nil
}
