package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type schedule struct {
	date  string
	para1 string
	para2 string
	para3 string
	para4 string
	para5 string

	kpara1 string
	kpara2 string
	kpara3 string
	kpara4 string
	kpara5 string
}
type changes struct {
	pars    []string
	prepods []string
	kabs    []string
}
type actually struct {
	para1 string
	para2 string
	para3 string
	para4 string
	para5 string
}

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
func trimToWord(s, word string) string {
	index := strings.Index(s, word)
	if index == -1 {
		return "" // слово не найдено
	}
	return s[index:]
}
func handleMainSchedule(url, linkText, fileName string) []string {
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

	tWeekDay := GetWeekday(0)
	for i, row := range rows {
		if strings.Contains(row[0], tWeekDay) {
			for j, col := range rows[0] {
				if strings.Contains(col, "3-ИС3") {
					for k := range 5 {
						Schedule = append(Schedule, rows[i+k][j])
					}

				}
			}
		}

	}
	return Schedule
}
func handleChangesSchedule(url, linkText, fileName string) []string {
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
	return Changes
}
func cleanString(s []string) []string {
	var result []string
	for _, part := range s {
		t := strings.ReplaceAll(part, "\n", " ")

		result = append(result, t)
	}
	return result

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
	date     string
	lessons  []string
	prepods  []string
	kabinets []string
}

func actuallyShedule() {
	sched := organizedChanges(handleMainSchedule(siteURL, "Расписание занятий на 1 семестр", "mainSchedule.xlsx"))
	chen := organizedChanges(handleChangesSchedule(siteURL, "Изменения в расписании", "changesSchedule.xlsx"))

	scheduleTable := actuallyTable{}

	for _, v := range sched {
		if strings.HasSuffix(v, ".") {
			scheduleTable.prepods = append(scheduleTable.prepods, v)
		} else {
			scheduleTable.lessons = append(scheduleTable.lessons, v)

		}
		scheduleTable.kabinets = append(scheduleTable.kabinets, "хз")

	}
	fmt.Println(sched)
	fmt.Println(chen)
	fmt.Println(scheduleTable)

	for i := 0; i < len(chen); i++ {
		item := chen[i]

		if strings.Contains(item, "ауд.") && i < len(scheduleTable.kabinets) {
			scheduleTable.kabinets[i-2] = item
			continue
		}
		if strings.Contains(item, "ауд.") && i > len(scheduleTable.kabinets) {
			scheduleTable.kabinets = append(scheduleTable.kabinets, item)
			continue
		}

		if strings.HasSuffix(item, ".") && i < len(scheduleTable.prepods) {
			scheduleTable.prepods[i-1] = item
			continue
		}
		if strings.HasSuffix(item, ".") && i > len(scheduleTable.prepods) {
			scheduleTable.prepods = append(scheduleTable.prepods, item)
			continue
		}
		if !strings.Contains(item, "ауд.") && !strings.HasSuffix(item, ".") && item != "-" && i < len(scheduleTable.lessons) {
			scheduleTable.lessons[i] = item
			continue
		}
		if !strings.Contains(item, "ауд.") && !strings.HasSuffix(item, ".") && item != "-" && i > len(scheduleTable.lessons) {
			scheduleTable.lessons = append(scheduleTable.lessons, item)
			continue
		}
		if item == "-" {
			continue
		}
	}
	fmt.Println(scheduleTable.lessons)
	fmt.Println(scheduleTable.prepods)
	fmt.Println(scheduleTable.kabinets)
}
