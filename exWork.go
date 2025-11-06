package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func GetTodayWeekday() string {
	weekdays := []string{
		"Воскресенье",
		"Понедельник",
		"Вторник",
		"Среда",
		"Четверг",
		"Пятница",
		"Суббота",
	}
	return weekdays[time.Now().Weekday()]
}

// FindRowByGroup ищет строку с groupName и убирает автоперенос в ячейках
func FindRowByGroup(filename, groupName string, sheetNum int) ([]string, error) {
	// Открываем файл Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	// Получаем список всех листов
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("в файле нет листов")
	}

	sheetName := sheetList[sheetNum]
	// Сначала получаем все координаты ячеек
	cols, _ := f.GetCols(sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строк листа %q: %w", sheetName, err)
	}
	var o []string
	// Поиск строки с groupName
	for i, r := range rows {
		if strings.Contains(r[0], GetTodayWeekday()) {
			for j, c := range cols {
				if strings.Contains(c[0], groupName) {
					for k := range 5 {
						if rows[i+k][j] == "" {
							break
						}
						o = append(o, rows[k+1][1], rows[i+k][j], "\n")

					}

				}
			}
		}

	}
	cleanedData := make([]string, len(o))

	for i, item := range o {
		// Удаляем все \n
		cleanedData[i] = strings.ReplaceAll(item, "\n", "  ")
		// Или удаляем \n только в конце: cleanedData[i] = strings.TrimSuffix(item, "\n")
	}

	return cleanedData, nil
	//return nil, fmt.Errorf("группа %q не найдена в файле", groupName)
}

func FindCollByGroup(filename, groupName string, sheetNum int) ([]string, error) {
	// Открываем файл Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	// Получаем список всех листов
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("в файле нет листов")
	}

	sheetName := sheetList[sheetNum]
	// Сначала получаем все координаты ячеек
	cols, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения столбов листа %q: %w", sheetName, err)
	}

	// Поиск строки с groupName
	for _, col := range cols {
		for _, cell := range col {
			if strings.EqualFold(cell, groupName) {
				return col, nil
			}

		}
	}

	return nil, fmt.Errorf("группа %q не найдена в файле", groupName)
}

func getChanges(name, linkText string) string {
	// 1. Находим ВТОРУЮ ссылку по тексту
	excelLink, err := FindDoc(siteURL, linkText)
	if err != nil {
		fmt.Printf("Ошибка поиска второй ссылки: %v\n", err)

	}
	fmt.Printf("Найденная вторая ссылка: %s\n", excelLink)

	// 2. Преобразуем URL Google Таблиц в URL экспорта
	exportURL, err := convertToGoogleExportURL(excelLink)
	if err != nil {
		fmt.Printf("Ошибка преобразования URL в формат экспорта: %v\n", err)

	}
	fmt.Printf("URL для экспорта: %s\n", exportURL)

	// 3. Скачиваем файл
	filePath, err := DownloadFile(exportURL, outputDir, "changesFile.xlsx")
	if err != nil {
		fmt.Printf("Ошибка скачивания файла: %v\n", err)

	}

	fmt.Printf("Файл успешно сохранён: %s\n", filePath)
	list, _ := FindRowByGroup("changesFile.xlsx", name, 0)
	fmt.Print(list)

	fmt.Println("")
	var outputList []string
	for i, l := range list {
		l = strings.ReplaceAll(l, "\r", " ")
		preparedValue := strings.ReplaceAll(l, "\n", " ")
		if preparedValue == "" {
			outputList = append(outputList, fmt.Sprintf("%d пара без изм.", i))
		} else {
			temp := fmt.Sprintf("%d пара: %s", i, preparedValue)
			outputList = append(outputList, temp)
		}
	}

	return strings.Join(outputList, "\n")

}
func trimToWord(s, word string) string {
	index := strings.Index(s, word)
	if index == -1 {
		return "" // слово не найдено
	}
	return s[index:]
}
