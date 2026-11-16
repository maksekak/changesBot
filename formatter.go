package main

import (
	"fmt"
	"strings"
	"time"
)

// weekEmoji возвращает эмодзи для типа недели
func weekEmoji(week string) string {
	switch week {
	case WeekBlue:
		return "🔵"
	case WeekRed:
		return "🔴"
	}
	return "⚪"
}

// FormatDaySchedule форматирует расписание на один день в виде компактного списка,
// который корректно отображается на экранах телефонов, с наложением актуальных изменений.
func FormatDaySchedule(groupName string, dayName string, date string, week string, lessons []Lesson, changesData ChangesData) string {
	var sb strings.Builder

	weekIcon := weekEmoji(week)
	weekLabel := ""
	switch week {
	case WeekBlue:
		weekLabel = "синяя"
	case WeekRed:
		weekLabel = "красная"
	}

	sb.WriteString(fmt.Sprintf("📅 Расписание на %s, %s (%s неделя %s)\n", dayName, date, weekLabel, weekIcon))
	sb.WriteString(fmt.Sprintf("👥 Группа: %s\n\n", groupName))

	// Проверяем, относятся ли изменения в таблице к целевой дате
	changesMatchDate := changesData.DateHeader != "" &&
		strings.Contains(strings.ToLower(changesData.DateHeader), strings.ToLower(date))

	changeMap := make(map[int]ChangeEntry)
	if changesMatchDate {
		for _, c := range changesData.Entries {
			changeMap[c.PairNumber] = c
		}
	}

	// Отфильтруем пары по текущей неделе
	var activeLessons []Lesson
	for _, lesson := range lessons {
		if lesson.Week != "" && lesson.Week != week {
			continue
		}
		activeLessons = append(activeLessons, lesson)
	}

	if len(activeLessons) == 0 && len(changeMap) == 0 {
		sb.WriteString("Пар нет.")
		if changesData.DateHeader != "" && !changesMatchDate {
			sb.WriteString(fmt.Sprintf("\n\nℹ️ В таблице изменений указана дата: %s.", changesData.DateHeader))
		}
		return sb.String()
	}

	pairIcons := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣"}
	appliedPairs := make(map[int]bool)

	for _, lesson := range activeLessons {
		icon := "•"
		if lesson.Number >= 1 && lesson.Number <= 6 {
			icon = pairIcons[lesson.Number-1]
		}

		timeStr := lesson.Time
		if timeStr == "" {
			timeStr = "—"
		}

		pairWeekIcon := "🔵🔴"
		if lesson.Week == WeekBlue {
			pairWeekIcon = "🔵"
		} else if lesson.Week == WeekRed {
			pairWeekIcon = "🔴"
		}

		line := fmt.Sprintf("%s %s — %s", icon, timeStr, lesson.Subject)
		if lesson.Teacher != "" {
			line += fmt.Sprintf(" — %s", lesson.Teacher)
		}
		if lesson.Room != "" {
			line += fmt.Sprintf(" — ауд. %s", lesson.Room)
		}

		if change, ok := changeMap[lesson.Number]; ok {
			appliedPairs[lesson.Number] = true
			if IsCancellation(change.Content) {
				line += " ❌ [ОТМЕНА]"
			} else {
				line += fmt.Sprintf(" 🔄 [ЗАМЕНА: %s]", change.Content)
			}
		}

		line += fmt.Sprintf(" [%s]", pairWeekIcon)
		sb.WriteString(line + "\n")
	}

	// Добавляем пары из таблицы изменений, которых не было в основном расписании на этот день
	for pairNum := 1; pairNum <= 6; pairNum++ {
		if appliedPairs[pairNum] {
			continue
		}
		change, ok := changeMap[pairNum]
		if !ok {
			continue
		}
		icon := pairIcons[pairNum-1]
		if IsCancellation(change.Content) {
			sb.WriteString(fmt.Sprintf("%s %d пара — ❌ [ОТМЕНА]\n", icon, pairNum))
		} else {
			sb.WriteString(fmt.Sprintf("%s %d пара — ➕ [ИЗМЕНЕНИЕ: %s]\n", icon, pairNum, change.Content))
		}
	}

	sb.WriteString("\n")
	if changesData.DateHeader == "" {
		sb.WriteString("ℹ️ Таблица изменений не загружена или пуста.")
	} else if !changesMatchDate {
		sb.WriteString(fmt.Sprintf("ℹ️ Изменений на %s пока нет (в таблице изменений дата: %s).", date, changesData.DateHeader))
	} else if len(changeMap) == 0 {
		sb.WriteString(fmt.Sprintf("✅ Изменений на %s для группы %s нет.", date, groupName))
	} else {
		sb.WriteString(fmt.Sprintf("ℹ️ Изменения на %s наложены на основное расписание.", date))
	}

	return sb.String()
}

// GetEffectiveScheduleDay возвращает день, дату и номер дня недели для отображения.
// До 17:00 — сегодня, после 17:00 — завтра.
func GetEffectiveScheduleDay() (dayName string, date string, weekday int) {
	now := time.Now()

	target := now
	if now.Hour() >= 17 {
		target = now.AddDate(0, 0, 1)
	}

	weekday = int(target.Weekday())

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

	date = fmt.Sprintf("%d %s", target.Day(), months[target.Month()-1])

	return
}

// FormatScheduleForTomorrow собирает полное сообщение с расписанием на сегодня/завтра
// (до 17:00 — сегодня, после 17:00 — завтра)
func FormatScheduleForTomorrow(groupName string, schedule GroupSchedule, changesData ChangesData, week string) string {
	_, date, weekday := GetEffectiveScheduleDay()

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

	return FormatDaySchedule(groupName, ruDay, date, week, lessons, changesData)
}