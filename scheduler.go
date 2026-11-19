package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartScheduler запускает фоновую проверку времени для рассылки в 7:00 и 20:00
func StartScheduler(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			now := time.Now()
			next := nextNotifyTime(now)
			duration := next.Sub(now)
			log.Printf("Следующая рассылка в %v (через %v)", next.Format("15:04 02.01.2006"), duration)

			time.Sleep(duration)

			sendDailySchedule(bot)
		}
	}()
}

// nextNotifyTime возвращает ближайшее время рассылки (7:00 или 20:00) относительно now.
func nextNotifyTime(now time.Time) time.Time {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	times := []time.Time{
		time.Date(now.Year(), now.Month(), now.Day(), NotifyHour, NotifyMinute, 0, 0, now.Location()),
		time.Date(now.Year(), now.Month(), now.Day(), NotifyHour2, NotifyMinute2, 0, 0, now.Location()),
	}

	var next time.Time
	for _, t := range times {
		if now.Before(t) {
			if next.IsZero() || t.Before(next) {
				next = t
			}
		}
	}

	if next.IsZero() {
		// Оба времени на сегодня уже прошли — следующее утро (7:00)
		next = today.AddDate(0, 0, 1).Add(time.Duration(NotifyHour2)*time.Hour + time.Duration(NotifyMinute2)*time.Minute)
	}

	return next
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