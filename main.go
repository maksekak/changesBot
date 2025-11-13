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
	switch text {
	case "/mainschedule":
		sched := organizedChanges(handleMainSchedule(siteURL, "Расписание занятий на 1 семестр", "mainSchedule.xlsx"))
		o := strings.Join(sched, "\n")
		msg := tgbotapi.NewMessage(message.Chat.ID, o)
		bot.Send(msg)
	case "/changes":
		chen := organizedChanges(handleChangesSchedule(siteURL, "Изменения в расписании", "changesSchedule.xlsx"))
		var out []string
		for _, o := range chen {
			if o == "-" || o == "ОТМЕНА" {
				o = o + "\n\n"
				out = append(out, o)
			} else {
				o = o + "\n"
				out = append(out, o)
			}
		}
		t := strings.Join(out, "")
		msg := tgbotapi.NewMessage(message.Chat.ID, t)
		bot.Send(msg)
	case "/actuallyschedule":
		act := actuallyShedule()
		ProcessStruct(&act)

		headers := []string{"  Пары  ", "  Преподы ", "  Кабинеты "}
		les := act.Lessons
		prep := act.Prepods
		kab := act.Kabinets

		var tableData [][]string
		for i := 0; i < len(prep); i++ {
			tableData = append(tableData, []string{
				les[i],
				fmt.Sprintf("%s", prep[i]),
				kab[i],
			})
		}
		table := buildMarkdownTable(headers, tableData)
		msg := tgbotapi.NewMessage(message.Chat.ID, table)
		msg.ParseMode = "Markdown"
		bot.Send(msg)

	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

}
