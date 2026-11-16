package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCommand обрабатывает команды бота
func handleCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	chatID := message.Chat.ID
	text := message.Text
	userName := message.From.FirstName

	log.Printf("%s (%d) wrote: %s", userName, chatID, text)

	args := strings.Fields(text)
	command := args[0]

	switch command {
	case "/start", "/start@collegeChangesBot":
		sendMessage(bot, chatID, `Привет! Я бот расписания колледжа.

Команды:
/setgroup <группа> — установить свою группу (например, /setgroup 4-NS1)
/mygroup — показать свою группу
/week blue|red — указать текущую неделю (дальше чередуется автоматически)
/notify — включить/выключить ежедневную рассылку в 20:00
/schedule — расписание на завтра с изменениями
/help — эта справка`)

	case "/help", "/help@collegeChangesBot":
		sendMessage(bot, chatID, `Команды бота:
/setgroup <группа> — установить свою группу
/mygroup — показать свою группу
/week blue|red — указать текущую неделю
/notify — вкл/выкл рассылку в 20:00
/schedule — расписание на завтра`)

	case "/setgroup", "/setgroup@collegeChangesBot":
		if len(args) < 2 {
			sendMessage(bot, chatID, "Укажите группу: /setgroup 4-NS1")
			return
		}
		groupName := args[1]
		if err := setGroup(chatID, groupName); err != nil {
			log.Printf("Ошибка сохранения группы: %v", err)
			sendMessage(bot, chatID, "Ошибка сохранения группы. Попробуйте позже.")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Группа %s сохранена.", groupName))

	case "/mygroup", "/mygroup@collegeChangesBot":
		groupName, ok := getGroup(chatID)
		if !ok || groupName == "" {
			sendMessage(bot, chatID, "Группа не установлена. Используйте /setgroup <группа>")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Ваша группа: %s", groupName))

	case "/week", "/week@collegeChangesBot":
		if len(args) < 2 {
			sendMessage(bot, chatID, "Укажите неделю: /week blue или /week red")
			return
		}
		week := strings.ToLower(args[1])
		if week != WeekBlue && week != WeekRed {
			sendMessage(bot, chatID, "Используйте /week blue или /week red")
			return
		}
		if err := setWeek(week); err != nil {
			log.Printf("Ошибка сохранения недели: %v", err)
			sendMessage(bot, chatID, "Ошибка сохранения недели. Попробуйте позже.")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Неделя установлена: %s. Далее неделя будет чередоваться автоматически каждые 7 дней.", week))

	case "/notify", "/notify@collegeChangesBot":
		enabled := toggleNotify(chatID)
		if enabled {
			sendMessage(bot, chatID, "Ежедневная рассылка включена. Буду присылать расписание на завтра в 20:00.")
		} else {
			sendMessage(bot, chatID, "Ежедневная рассылка выключена.")
		}

	case "/schedule", "/schedule@collegeChangesBot", "/завтра", "/завтра@collegeChangesBot":
		handleScheduleRequest(bot, chatID)

	default:
		// Неизвестная команда — игнорируем
	}
}

func handleScheduleRequest(bot *tgbotapi.BotAPI, chatID int64) {
	groupName, ok := getGroup(chatID)
	if !ok || groupName == "" {
		sendMessage(bot, chatID, "Группа не установлена. Используйте /setgroup <группа>")
		return
	}

	week, err := getCurrentWeek()
	if err != nil {
		sendMessage(bot, chatID, err.Error())
		return
	}

	schedule, err := ParseMainSchedule(groupName)
	if err != nil {
		log.Printf("Ошибка парсинга расписания: %v", err)
		sendMessage(bot, chatID, "Ошибка загрузки расписания. Попробуйте позже.")
		return
	}

	_, _, _ = GetTomorrowInfo()

	// Для изменений нам нужно знать номер пары, но в таблице изменений нет привязки к дню недели
	// Поэтому загружаем все изменения для группы
	changes, err := ParseChanges(groupName, "")
	if err != nil {
		log.Printf("Ошибка парсинга изменений: %v", err)
		// Не критично, продолжаем без изменений
		changes = nil
	}

	msg := FormatScheduleForTomorrow(groupName, schedule, changes, week)
	sendMessage(bot, chatID, msg)
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки в чат %d: %v", chatID, err)
	}
}

// HandleUpdate — основная точка входа для обработки обновлений
func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}

	message := update.Message
	if message.Text == "" {
		return
	}

	handleCommand(message, bot)
}