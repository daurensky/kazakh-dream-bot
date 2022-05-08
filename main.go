package main

import (
	"database/sql"
	"fmt"
	"github.com/daurensky/kazakh-dream-bot/db"
	"github.com/daurensky/kazakh-dream-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var messages = map[string]string{
	"dont_understand":       "Я вас не понял",
	"order_details":         "Заполните информацию о доставке в формате: /order_details Адрес, Ваш номер телефона, Ваше ФИО",
	"order_details_updated": "Данные о доставке обновлены",
	"order_details_created": "Данные о доставке созданы, можно продолжить заказывать",
	"order_created":         "Заказ создан. Посмотреть заказы можно с помощью /orders",
	"cart_empty":            "Корзина пуста. Пора её наполнить /menu",
	"orders_empty":          "Вы пока ничего не заказывали",
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))

	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "menu":
				products, err := db.GetProducts()

				if err != nil {
					panic(err)
				}

				for _, product := range products {
					file := tgbotapi.FileURL(product.PhotoUrl)

					photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, file)

					photoMsg.Caption = fmt.Sprintf(
						"*Название*: %s\n"+
							"*Цена*: %.2f тг.\n"+
							"*Состав*: %s",
						product.Name, product.Price, strings.Join(product.Composition, ", "))

					photoMsg.ParseMode = "markdown"

					addToCartButton := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Добавить в корзину", "to_cart:"+strconv.Itoa(product.Id)),
						),
					)

					photoMsg.ReplyMarkup = addToCartButton

					if _, err := bot.Send(photoMsg); err != nil {
						panic(err)
					}
				}
			case "cart":
				products, err := db.GetClientProducts(update.Message.From.ID, false)

				if err != nil {
					panic(err)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["cart_empty"])

				if len(products) != 0 {
					var rows []string

					for _, product := range products {
						rows = append(rows, fmt.Sprintf("%s за %.2f тг.", product.Name, product.Price))
					}

					msg.Text = strings.Join(rows, "\n")

					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"Создать заказ",
								"do_order:"+strconv.Itoa(int(update.Message.From.ID)),
							),
						),
					)
				}

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "order_details":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["order_details"])

				clientInfo := strings.Split(update.Message.Text, ",")

				if len(clientInfo) == 3 {
					address := strings.TrimSpace(strings.Replace(clientInfo[0], "/order_details", "", 1))
					phone := strings.TrimSpace(clientInfo[1])
					name := strings.TrimSpace(clientInfo[2])

					if address != "" && phone != "" && name != "" {
						client := models.Client{
							TelegramId: update.Message.From.ID,
							Name:       name,
							Phone:      phone,
							Address:    address,
						}

						_, err := db.ShowClient(update.Message.From.ID)

						if err == sql.ErrNoRows {
							err := db.StoreClient(client)

							if err != nil {
								panic(err)
							}

							msg.Text = messages["order_details_created"]

							msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData(
										"Создать заказ",
										"do_order:"+strconv.Itoa(int(update.Message.From.ID)),
									),
								),
							)
						} else if err == nil {
							err = db.UpdateClient(client)

							if err != nil {
								panic(err)
							}

							msg.Text = messages["order_details_updated"]
						} else {
							panic(err)
						}
					}
				}

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "orders":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["orders_empty"])

				client, err := db.ShowClient(update.Message.From.ID)

				if err == nil {

					orders, err := db.GetClientOrders(update.Message.From.ID)

					if err != nil {
						panic(err)
					}

					if len(orders) != 0 {
						var rows []string

						for _, order := range orders {
							var composition []string

							for _, product := range order.Products {
								composition = append(composition, product.Name)
							}

							createdAt, err := order.CreatedAtText()

							if err != nil {
								panic(err)
							}

							rows = append(rows, fmt.Sprintf(
								"_%s_: *%s*. %s, %s, %s, %s",
								order.StatusText(),
								strings.Join(composition, ", "),
								client.Name,
								client.Phone,
								client.Address,
								createdAt,
							))
						}

						msg.Text = strings.Join(rows, "\n")
						msg.ParseMode = "markdown"
					}
				} else if err != nil && err != sql.ErrNoRows {
					panic(err)
				}

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages["dont_understand"])

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		} else if update.CallbackQuery != nil {
			request := strings.Split(update.CallbackQuery.Data, ":")

			switch request[0] {
			case "to_cart":
				msg := "Заказ добавлен в корзину!"

				productId, err := strconv.Atoi(request[1])

				if err != nil {
					panic(err)
				}

				product, err := db.ShowProduct(productId)

				if err != nil {
					panic(err)
				}

				err = db.StoreCart(update.CallbackQuery.From.ID, product)

				if err != nil {
					panic(err)
				}

				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, msg)

				if _, err := bot.Request(callback); err != nil {
					panic(err)
				}
			case "do_order":
				var text string

				client, err := db.ShowClient(update.CallbackQuery.From.ID)

				if err == sql.ErrNoRows {
					text = messages["order_details"]
				} else if err == nil {
					products, err := db.GetClientProducts(client.TelegramId, true)

					if len(products) == 0 {
						text = messages["cart_empty"]
					} else if err == nil {
						orderId, err := db.StoreOrder(client)

						if err != nil {
							panic(err)
						}

						err = db.StoreOrderProducts(orderId, products)

						if err != nil {
							panic(err)
						}

						err = db.DestroyCart(client.TelegramId)

						if err != nil {
							panic(err)
						}

						text = messages["order_created"]
					} else {
						panic(err)
					}
				} else {
					panic(err)
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}
}
