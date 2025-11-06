package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err = tgbotapi.NewBotAPI(os.Getenv("token"))
	// если ошибка инициализации паникуем
	if err != nil {
		log.Panic(err)
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	for {

		handleUpdate(<-updates)
	}

}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)

	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}
	if text == "/changes" || text == "/changes@collegeChangesBot" {
		temp := getChanges("3-ИС3")
		fmt.Println(temp)
		if temp == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Изменений нет")
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, trimToWord(temp, "1 пара"))
			bot.Send(msg)
		}

	}
	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

}
