package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	// APIToken for Telegram YourFarmBot
	APIToken = "194047760:AAF73o66eOYH98GQoTTtZDd4-dkefjw663I"
	// HeartbeatImage is a picture of a heartbeat
	HeartbeatImage = "./image/heartbeats.jpg"
	// PiggyImage is a picture of a pig
	PiggyImage = "./image/piggy.jpg"
	zipcode    = "78705"
)

type marketResponse struct {
	MapsLink string `json:"GoogleLink"`
	Address  string `json:"Address"`
	Schedule string `json:"Schedule"`
	Products string `json:"Products"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(APIToken)
	if err != nil {
		log.Printf("Error: Could not create bot: %v\n", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 15

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Error: Could not get update chan: %v\n", err)
	}

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		log.Printf("\tMESSAGE: %s\n\t\t\tID: %d\n\t\t\tCaption: %s\n\t\t\tDate: %d\n", update.Message.Text, update.Message.MessageID, update.Message.Caption, update.Message.Date)

		switch update.Message.Text {
		case "hello":
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "goodbye!")
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {
					log.Printf("Error: Could not send message: %v\n", err)
				}

			}
		case "hello?":
			{
				photoUpload := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, HeartbeatImage)
				photoUpload.Caption = "It's alive!!"
				photoUpload.FileID = "heartbeat"
				if msg, err := bot.Send(photoUpload); err != nil {
					log.Println(err)
				} else {
					log.Println(msg)
				}
			}

		case "piggy":
			{
				photoUpload := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, PiggyImage)
				photoUpload.Caption = "Little piggy!"
				photoUpload.FileID = "piggy"
				if msg, err := bot.Send(photoUpload); err != nil {
					log.Println(err)
				} else {
					log.Println(msg)
				}

			}

		case "nearest market":
			{
				//TODO: location: "http://search.ams.usda.gov/farmersmarkets/v1/data.svc/locSearch?lat=" + lat + "&lng=" + lng
				resp, err := http.Get("http://search.ams.usda.gov/farmersmarkets/v1/data.svc/zipSearch?zip=" + zipcode)
				if err != nil {
					log.Printf("Could net search zipcode %s: %v", zipcode, err)
				}
				defer func() {
					if err := resp.Body.Close(); err != nil {
						log.Println(err)
					}
				}()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
				}

				newMarket := &marketResponse{}
				if err := json.Unmarshal(body, newMarket); err != nil {
					log.Println(err)
				}

				log.Println("response: " + newMarket.Address)
			}
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nobody here...")
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error: Could not send message: %v\n", err)
			}
		}
	}
}
