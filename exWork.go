package main

import (
	"fmt"
	"reflect"

	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func GetWeekday(dayChange int) string {
	weekdays := []string{
		"Воскресенье",
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}
	// Получаем текущий день недели (0 = воскресенье, 6 = суббота)
	today := int(time.Now().Weekday())
	// Вычисляем индекс завтрашнего дня (с учётом перехода через неделю)
	tomorrow := (today + dayChange) % 7

	return weekdays[tomorrow]
}
func getDate(p int) string {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, p)

	day := tomorrow.Day()
	month := tomorrow.Month()

	// Массив с названиями месяцев в нижнем регистре
	months := []string{
		"январь", "февраль", "март", "апрель", "май", "июнь",
		"июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь",
	}

	// Получаем название месяца по номеру (Month() возвращает 1–12)
	monthName := months[int(month)-1]

	// Форматируем день с ведущим нулём (двузначное число)
	return fmt.Sprintf("%d %s", day, monthName)
}

func handleMainSchedule(url, linkText, fileName, day string) []string {
	GetSchedule(url, linkText, fileName, 3)
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	sheets := f.GetSheetList()
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		fmt.Print(err)
	}
	var Schedule []string

	tWeekDay := day

	for i, row := range rows {
		if strings.Contains(row[0], tWeekDay) {
			for j, col := range rows[0] {
				if strings.Contains(col, "3-ИС3") {
					for k := range 4 {
						Schedule = append(Schedule, rows[i+k][j])
					}

				}
			}
		}

	}
	return Schedule
}
func handleChangesSchedule(url, linkText, fileName string) ([]string, string) {
	GetSchedule(url, linkText, fileName, 0)
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	sheets := f.GetSheetList()
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		fmt.Print(err)
	}
	var temp []string
	var Changes []string
	for i, row := range rows {
		if strings.Contains(row[0], "3-ИС3") {
			for _, col := range rows[i] {
				if col == "" {
					temp = append(temp, "-")
					continue
				}
				temp = append(temp, col)
				Changes = temp[1:]
			}

		}
	}
	day := extractDayOfWeek(rows[0][0])

	return Changes, day
}

func organizedChanges(b []string) []string {
	var act []string
	var d []string
	for _, c := range b {

		d = strings.Split(c, "\n")
		act = append(act, d...)
	}
	for i, v := range act {
		if v == "Иност. язык" {
			act[i] = act[i] + " " + act[i+1]
			act = append(act[:i+1], act[i+2:]...)
		}
	}
	return act
}

type actuallyTable struct {
	Lessons  []string
	Prepods  []string
	Kabinets []string
}

func extractDayOfWeek(text string) string {
	// Список дней недели на русском (в нижнем регистре для сравнения)
	days := []string{
		"понедельник", "вторник", "среда", "четверг",
		"пятница", "суббота", "воскресенье",
	}

	textLower := strings.ToLower(text) // приводим к нижнему регистру

	for _, day := range days {
		if strings.Contains(textLower, day) {
			return strings.Title(day)
		}
	}
	return "" // не найдено
}
func actuallyShedule() (actuallyTable, interface{}) {
	scheduleTable := actuallyTable{}
	c, day := handleChangesSchedule(siteURL, "Изменения в расписании", "changesSchedule.xlsx")
	if day == "" {
		return scheduleTable, "-"
	}
	chen := organizedChanges(c)
	sched := organizedChanges(handleMainSchedule(siteURL, "Расписание занятий на 1 семестр", "mainSchedule.xlsx", day))

	for _, v := range sched {
		if strings.HasSuffix(v, ".") {
			scheduleTable.Prepods = append(scheduleTable.Prepods, v)
		} else {
			scheduleTable.Lessons = append(scheduleTable.Lessons, v)

		}
		scheduleTable.Kabinets = append(scheduleTable.Kabinets, "хз")

	}
	l, p, k := 0, 0, 0
	for i := 0; i < len(chen); i++ {
		item := chen[i]

		if strings.Contains(item, "ауд.") {
			scheduleTable.Kabinets[k] = item
			k++
			continue
		}

		if strings.HasSuffix(item, ".") {
			scheduleTable.Prepods[p] = item
			p++
			continue
		}

		if !strings.Contains(item, "ауд.") && !strings.HasSuffix(item, ".") && item != "-" {
			scheduleTable.Lessons[l] = item
			l++
			continue
		}

		if item == "-" {
			k++
			p++
			l++
			continue
		}
	}

	scheduleTable.Kabinets = scheduleTable.Kabinets[:len(scheduleTable.Lessons)]
	return scheduleTable, nil
}
func buildMarkdownTable(headers []string, data [][]string) string {
	var lines []string

	// Шаг 1: Определяем максимальную ширину для каждого столбца
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len([]rune(escapeMarkdown(header)))
	}

	for _, row := range data {
		for i, cell := range row {
			cellWidth := len([]rune(escapeMarkdown(cell)))
			if cellWidth > colWidths[i] {
				colWidths[i] = cellWidth
			}
		}
	}

	// Шаг 2: Формируем строку заголовков
	headerLine := "|"
	for i, header := range headers {
		escapedHeader := escapeMarkdown(header)
		paddedHeader := padToWidth(escapedHeader, colWidths[i])
		headerLine += "" + paddedHeader + "|"
	}
	lines = append(lines, headerLine)

	// Шаг 3: Формируем разделительную строку
	sepLine := "|"
	for width := range colWidths {
		sepLine += strings.Repeat("-", colWidths[width]) + "|" // +2 для пробелов по краям
	}
	lines = append(lines, sepLine)

	// Шаг 4: Формируем строки данных
	for _, row := range data {
		rowLine := "|"
		for i, cell := range row {
			escapedCell := escapeMarkdown(cell)
			paddedCell := padToWidth(escapedCell, colWidths[i])
			rowLine += "" + paddedCell + "|"
		}
		lines = append(lines, rowLine)
	}

	return fmt.Sprintf("```Актуальная-таблица\n%s\n```", strings.Join(lines, "\n"))
}
func padToWidth(s string, width int) string {
	runeCount := len([]rune(s))
	if runeCount >= width {
		return s
	}
	spacesNeeded := width - runeCount
	return s + strings.Repeat(" ", spacesNeeded)
}
func escapeMarkdown(text string) string {
	return strings.ReplaceAll(text, "|", "\\|")
}
func TrimToFirstWord(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return words[0]
	}

	// Проверяем второе слово: содержит ли цифры или оканчивается на точку
	secondWord := words[1]
	hasDigits := strings.ContainsAny(secondWord, "0123456789")
	endsWithDot := len(secondWord) > 0 && secondWord[len(secondWord)-1] == '.'

	if hasDigits || endsWithDot {
		// Сохраняем первые два слова
		return words[0] + " " + secondWord
	}

	// Иначе оставляем только первое слово
	return words[0]
}

// processStruct обходит структуру и обрабатывает все слайсы строк []string
func ProcessStruct(v interface{}) {
	rv := reflect.ValueOf(v)

	// Проверяем, что передан указатель на структуру
	if rv.Kind() != reflect.Ptr {
		panic("Expected a pointer")
	}
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		panic("Expected a struct")
	}

	// Проходим по всем полям структуры
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)

		// Обрабатываем только слайсы строк []string
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
			// Получаем текущий слайс
			slice := field
			for j := 0; j < slice.Len(); j++ {
				// Берем строку, обрезаем, записываем обратно
				str := slice.Index(j).String()
				trimmed := TrimToFirstWord(str)
				slice.Index(j).SetString(trimmed)
			}
		}
	}
}
