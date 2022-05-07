package db

func DestroyCart(telegramId int64) error {
	db, err := Connect()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM kazakh_dream.public.cart WHERE telegram_id = $1", telegramId)

	if err != nil {
		return err
	}

	return nil
}
