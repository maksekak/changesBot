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
	command := strings.ToLower(args[0])
	if idx := strings.Index(command, "@"); idx != -1 {
		botMention := command[idx+1:]
		if bot.Self.UserName != "" && !strings.EqualFold(botMention, bot.Self.UserName) {
			return
		}
		command = command[:idx]
	}

	switch command {
	case "/start":
		sendMessage(bot, chatID, `Привет! Я бот расписания колледжа.

Команды:
/setgroup <группа> — установить группу для этого чата (например, /setgroup 4-ИС3)
/mygroup — показать текущую группу
/groups — список доступных групп
/weekset ДД.ММ blue|red — установить дату и цвет недели (например, /weekset 01.09 blue)
/notify — включить/выключить ежедневную рассылку на завтра в 20:00
/schedule — расписание с изменениями (до 17:00 на сегодня, после 17:00 на завтра)
/help — эта справка`)

	case "/help":
		sendMessage(bot, chatID, `Команды бота:
/setgroup <группа> — установить группу (например, /setgroup 4-ИС3)
/mygroup — показать текущую группу
/groups — список доступных групп
/weekset ДД.ММ blue|red — установить дату и цвет недели
/notify — вкл/выкл рассылку в 20:00
/schedule — расписание с изменениями (до 17:00 сегодня, после 17:00 завтра)`)

	case "/groups":
		sendMessage(bot, chatID, "Доступные группы: "+strings.Join(GetAvailableGroups(), ", "))

	case "/setgroup":
		if len(args) < 2 {
			sendMessage(bot, chatID, "Укажите группу, например: /setgroup 4-ИС3\nДоступные группы: "+strings.Join(GetAvailableGroups(), ", "))
			return
		}
		schedule, err := GetHardcodedSchedule(args[1])
		if err != nil {
			sendMessage(bot, chatID, err.Error())
			return
		}
		if err := setGroup(chatID, schedule.Group); err != nil {
			log.Printf("Ошибка сохранения группы: %v", err)
			sendMessage(bot, chatID, "Ошибка сохранения группы. Попробуйте позже.")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Группа %s сохранена.", schedule.Group))

	case "/mygroup":
		groupName, ok := getGroup(chatID)
		if !ok || groupName == "" {
			sendMessage(bot, chatID, "Группа не установлена. Используйте /setgroup <группа>")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Текущая группа: %s", groupName))

	case "/week":
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
		sendMessage(bot, chatID, fmt.Sprintf("Неделя установлена: %s. Далее неделя будет чередоваться автоматически каждый понедельник.", week))

	case "/weekset":
		if len(args) < 3 {
			sendMessage(bot, chatID, "Укажите дату и неделю: /weekset 01.09 blue")
			return
		}
		dateStr := args[1]
		week := strings.ToLower(args[2])
		if week != WeekBlue && week != WeekRed {
			sendMessage(bot, chatID, "Используйте /weekset ДД.ММ blue|red, например /weekset 01.09 blue")
			return
		}
		refDate, err := parseDateDDMM(dateStr)
		if err != nil {
			sendMessage(bot, chatID, err.Error())
			return
		}
		if err := setWeekWithDate(refDate, week); err != nil {
			log.Printf("Ошибка сохранения недели: %v", err)
			sendMessage(bot, chatID, "Ошибка сохранения недели. Попробуйте позже.")
			return
		}
		sendMessage(bot, chatID, fmt.Sprintf("Неделя установлена: %s с даты %s. Далее неделя будет чередоваться автоматически каждый понедельник.", week, dateStr))

	case "/notify":
		enabled := toggleNotify(chatID)
		if enabled {
			sendMessage(bot, chatID, "Ежедневная рассылка включена. Буду присылать расписание на завтра в 20:00.")
		} else {
			sendMessage(bot, chatID, "Ежедневная рассылка выключена.")
		}

	case "/schedule", "/завтра":
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
		log.Printf("Ошибка получения расписания: %v", err)
		sendMessage(bot, chatID, err.Error())
		return
	}

	changesData, err := ParseChanges(schedule.Group, "")
	if err != nil {
		log.Printf("Ошибка парсинга изменений: %v", err)
	}

	msg := FormatScheduleForTomorrow(schedule.Group, schedule, changesData, week)
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