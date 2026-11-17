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

// FormatDaySchedule форматирует расписание на один день в новом виде:
//
//	Среда, 18 ноября (красная неделя)
//	4-ИС3
//
//	2️⃣ МДК 01.04 (10:15 - 11:50)
//	   Евстигнеев А.С.
//
//	0. Семьеведение (12:20 - 12:35)
//
//	3️⃣ МДК 11.01 (12:45 - 14:20)
//	   Бобков Д.И.
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

	sb.WriteString(fmt.Sprintf("%s, %s (%s неделя %s)\n", dayName, date, weekLabel, weekIcon))
	sb.WriteString(fmt.Sprintf("%s\n\n", groupName))

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

	for i, lesson := range activeLessons {
		if i > 0 {
			sb.WriteString("\n")
		}

		icon := "•"
		if lesson.Number >= 1 && lesson.Number <= 6 {
			icon = pairIcons[lesson.Number-1]
		}

		timeStr := lesson.Time
		if timeStr == "" {
			timeStr = "—"
		}

		sb.WriteString(fmt.Sprintf("%s %s (%s)\n", icon, lesson.Subject, timeStr))

		teacherLine := "   "
		if lesson.Teacher != "" {
			teacherLine += lesson.Teacher
		}
		if lesson.Room != "" {
			if lesson.Teacher != "" {
				teacherLine += " | ауд. " + lesson.Room
			} else {
				teacherLine += "ауд. " + lesson.Room
			}
		}
		sb.WriteString(teacherLine + "\n")

		if change, ok := changeMap[lesson.Number]; ok {
			appliedPairs[lesson.Number] = true
			if IsCancellation(change.Content) {
				sb.WriteString("   ❌ Отмена\n")
			} else {
				sb.WriteString(fmt.Sprintf("   🔄 Замена: %s\n", change.Content))
			}
		}
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
		sb.WriteString("\n")
		icon := pairIcons[pairNum-1]
		if IsCancellation(change.Content) {
			sb.WriteString(fmt.Sprintf("%s %d пара — ❌ Отмена\n", icon, pairNum))
		} else {
			sb.WriteString(fmt.Sprintf("%s %d пара — 🔄 Замена: %s\n", icon, pairNum, change.Content))
		}
	}

	if changesData.DateHeader == "" {
		sb.WriteString("\nℹ️ Таблица изменений не загружена или пуста.")
	} else if !changesMatchDate {
		sb.WriteString(fmt.Sprintf("\nℹ️ Изменений на %s пока нет (в таблице изменений дата: %s).", date, changesData.DateHeader))
	} else if len(changeMap) == 0 {
		sb.WriteString(fmt.Sprintf("\n✅ Изменений на %s для группы %s нет.", date, groupName))
	} else {
		sb.WriteString(fmt.Sprintf("\nℹ️ Изменения на %s наложены на основное расписание.", date))
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