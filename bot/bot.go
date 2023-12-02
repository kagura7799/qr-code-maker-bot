package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/url"
)

const TOKEN = "6821853023:AAFDhL3b8cEWKbo7sgkLZVYP8lLEZGAa-nQ"

func main() {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(bot, update.Message)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Привет, "+msg.From.FirstName+"! Добро пожаловать в бота, который создаст QR-код для твоей ссылки."))
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Просто отправь мне свою ссылку и в ответ ты получишь QR-Code на свой URL."))

	default:
		handleUserURL(bot, msg)
	}
}

func handleUserURL(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userURL := msg.Text

	if userURL != "/start" && checkURL(userURL) {
		jpegData, err := createQRCode(userURL)
		if err != nil {
			log.Println(err)
			return
		}

		photoConfig := tgbotapi.PhotoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID: msg.Chat.ID,
				},
				File: tgbotapi.FileBytes{
					Name:  "qr.jpg",
					Bytes: jpegData,
				},
			},
			Caption: "Ваш QR-Code по запросу.",
		}

		_, err = bot.Send(photoConfig)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Неправильный URL, попробуй еще раз"))
	}
}

func checkURL(input string) bool {
	fmt.Println("checkURL: ", input)
	u, err := url.Parse(input)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func createQRCode(userURL string) ([]byte, error) {
	q, err := qrcode.New(userURL, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	img := q.Image(256)
	img = resize.Resize(256, 256, img, resize.Lanczos3)

	rgba := image.NewRGBA(img.Bounds())
	drawImage := rgba.SubImage(img.Bounds()).(*image.RGBA)
	draw.Draw(drawImage, drawImage.Bounds(), img, image.Point{}, draw.Over)

	var jpegData []byte
	buffer := &bytes.Buffer{}
	err = jpeg.Encode(buffer, rgba, nil)
	if err != nil {
		return nil, err
	}
	jpegData = buffer.Bytes()

	return jpegData, nil
}
