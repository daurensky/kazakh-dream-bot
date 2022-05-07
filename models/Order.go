package models

import "time"

type Order struct {
	Id        int
	Status    string
	CreatedAt string
	ClientId  int
	Products  []Product
}

func (o Order) StatusText() string {
	switch o.Status {
	case "PREPARING":
		return "Готовится"
	case "SENT":
		return "Отправлено курьером"
	case "DELIVERED":
		return "Доставлено"
	default:
		return "Неизвестный статус заказа"
	}
}

func (o Order) CreatedAtText() (string, error) {
	t, err := time.Parse(time.RFC3339, o.CreatedAt)

	if err != nil {
		return "", err
	}

	return t.Format("02.01.06 03:04"), nil
}
