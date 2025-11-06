package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	bot     *tgbotapi.BotAPI
	siteURL = "https://tgiek.ru/studentam" // URL страницы с ссылкой

	outputDir = "" // Директория для сохранения файла
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
		temp := getChanges("3-ИС3", "Изменения в расписании")
		fmt.Println(temp)
		if temp == "" {
			msgg := tgbotapi.NewMessage(message.Chat.ID, "Изменений нет")
			bot.Send(msgg)
			c, _ := FindRowByGroup("scheduleFile.xlsx", "3-ИС3", 0)
			t := strings.Join(c, "\n")
			msg := tgbotapi.NewMessage(message.Chat.ID, t)
			bot.Send(msg)
		} else {
			msgg := tgbotapi.NewMessage(message.Chat.ID, "Изменения")
			bot.Send(msgg)
			msg := tgbotapi.NewMessage(message.Chat.ID, trimToWord(temp, "1 пара"))
			bot.Send(msg)
		}

	}
	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

}
