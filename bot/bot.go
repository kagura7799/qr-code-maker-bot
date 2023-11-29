package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/url"
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
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Просто отправь мне свою ссылку и в ответ ты получишь QR-Code на свой URL."))

			for {
				userURLMessage := <-updates
				userURL := userURLMessage.Message.Text
				fmt.Println(userURL)

				if checkURL(userURL) {
					fmt.Println("true url")

					jpegData, err := createQRCode(userURL)
					if err != nil {
						log.Fatal(err)
					}

					photoConfig := tgbotapi.PhotoConfig{
						BaseFile: tgbotapi.BaseFile{
							BaseChat: tgbotapi.BaseChat{
								ChatID: update.Message.Chat.ID,
							},
							File: tgbotapi.FileBytes{
								Name:  "qr.jpg",
								Bytes: jpegData,
							},
						},
						Caption: "Ваш QR-Code по запросу.", // Опциональный заголовок
					}

					_, err = bot.Send(photoConfig)
					if err != nil {
						log.Fatal(err)
					}

				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный URL, попробуй еще раз"))
				}

			}
		}
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
