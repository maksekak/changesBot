package main

import (
	"fmt"
	"strings"
	"time"
)

// FormatDaySchedule форматирует расписание на один день для вывода
func FormatDaySchedule(groupName string, dayName string, date string, week string, lessons []Lesson, changes []ChangeEntry) string {
	var sb strings.Builder

	weekLabel := ""
	switch week {
	case WeekBlue:
		weekLabel = "синяя"
	case WeekRed:
		weekLabel = "красная"
	}

	sb.WriteString(fmt.Sprintf("Расписание на %s, %s (%s неделя)\n", dayName, date, weekLabel))
	sb.WriteString(fmt.Sprintf("Группа: %s\n\n", groupName))

	if len(lessons) == 0 {
		sb.WriteString("Пар нет.")
		return sb.String()
	}

	// Строим карту изменений по номеру пары
	changeMap := make(map[int]ChangeEntry)
	for _, c := range changes {
		changeMap[c.PairNumber] = c
	}

	for _, lesson := range lessons {
		// Фильтруем по неделе
		if lesson.Week != "" && lesson.Week != week {
			continue
		}

		sb.WriteString(fmt.Sprintf("%d. %s", lesson.Number, lesson.Subject))
		if lesson.Teacher != "" {
			sb.WriteString(fmt.Sprintf(" — %s", lesson.Teacher))
		}
		if lesson.Room != "" {
			sb.WriteString(fmt.Sprintf(" — ауд. %s", lesson.Room))
		}

		// Проверяем изменения
		if change, ok := changeMap[lesson.Number]; ok {
			if IsCancellation(change.Content) {
				sb.WriteString(" [ОТМЕНА]")
			} else {
				sb.WriteString(fmt.Sprintf(" [ЗАМЕНА: %s]", change.Content))
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// GetTomorrowInfo возвращает название дня, дату и номер дня недели для завтра
func GetTomorrowInfo() (dayName string, date string, weekday int) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	weekday = int(tomorrow.Weekday())

	days := []string{
		"Воскресенье",
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}

	dayName = days[weekday]

	months := []string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}

	date = fmt.Sprintf("%d %s", tomorrow.Day(), months[tomorrow.Month()-1])

	return
}

// FormatScheduleForTomorrow собирает полное сообщение с расписанием на завтра
func FormatScheduleForTomorrow(groupName string, schedule GroupSchedule, changes []ChangeEntry, week string) string {
	_, date, weekday := GetTomorrowInfo()

	// Превращаем номер дня недели в название (понедельник = 1, воскресенье = 0)
	dayNames := map[int]string{
		1: "Понедельник",
		2: "Вторник",
		3: "Среда",
		4: "Четверг",
		5: "Пятница",
		6: "Суббота",
		0: "Воскресенье",
	}

	ruDay := dayNames[weekday]
	lessons := schedule.Days[ruDay]

	return FormatDaySchedule(groupName, ruDay, date, week, lessons, changes)
}