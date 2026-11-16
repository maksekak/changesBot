package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartScheduler запускает фоновую проверку времени для рассылки в 20:00
func StartScheduler(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			now := time.Now()
			// Вычисляем время до следующей 20:00
			next := time.Date(now.Year(), now.Month(), now.Day(), NotifyHour, NotifyMinute, 0, 0, now.Location())
			if now.After(next) {
				next = next.AddDate(0, 0, 1)
			}

			duration := next.Sub(now)
			log.Printf("Следующая рассылка в %v (через %v)", next.Format("15:04 02.01.2006"), duration)

			time.Sleep(duration)

			sendDailySchedule(bot)
		}
	}()
}

// sendDailySchedule отправляет расписание на завтра всем чатам с включёнными уведомлениями
func sendDailySchedule(bot *tgbotapi.BotAPI) {
	groups := loadGroups()

	week, err := getCurrentWeek()
	if err != nil {
		log.Printf("Не удалось определить неделю: %v", err)
		return
	}

	for chatID, settings := range groups {
		if !settings.Notify || settings.GroupName == "" {
			continue
		}

		schedule, err := ParseMainSchedule(settings.GroupName)
		if err != nil {
			log.Printf("Ошибка парсинга расписания для чата %d: %v", chatID, err)
			continue
		}

		changesData, err := ParseChanges(schedule.Group, "")
		if err != nil {
			log.Printf("Ошибка парсинга изменений для чата %d: %v", chatID, err)
		}

		msg := FormatScheduleForTomorrow(schedule.Group, schedule, changesData, week)
		sendTelegramMessage(bot, chatID, msg)
	}
}

// sendTelegramMessage отправляет сообщение через бота
func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки в чат %d: %v", chatID, err)
	}
}