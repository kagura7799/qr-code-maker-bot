package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	_ "os"
)

const TOKEN = "6821853023:AAFDhL3b8cEWKbo7sgkLZVYP8lLEZGAa-nQ"

func main() {
	var err error
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.FirstName+"! Добро пожаловать в бота который создаст QR-код для твоей ссылки."))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери действие:")
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Сгенерировать QR-Code"),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)

		case "test":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "да тест работает"))

		}
	}
}

func createQRCode() {

}
