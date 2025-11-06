package main

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// FindRowByGroup ищет строку с groupName и убирает автоперенос в ячейках
func FindRowByGroup(filename, groupName string) ([]string, error) {
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

	sheetName := sheetList[0]
	// Сначала получаем все координаты ячеек
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строк листа %q: %w", sheetName, err)
	}

	// Поиск строки с groupName
	for _, row := range rows {
		for _, cell := range row {
			if strings.EqualFold(cell, groupName) {
				return row, nil
			}

		}
	}

	return nil, fmt.Errorf("группа %q не найдена в файле", groupName)
}
func getChanges(name string) string {
	siteURL := "https://tgiek.ru/studentam" // URL страницы с ссылкой
	linkText := "Изменения в расписании"    // Текст ссылки, по которому ищем
	outputDir := ""                         // Директория для сохранения файла

	// 1. Находим ВТОРУЮ ссылку по тексту
	excelLink, err := FindSecondLinkByText(siteURL, linkText)
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
	filePath, err := DownloadFile(exportURL, outputDir)
	if err != nil {
		fmt.Printf("Ошибка скачивания файла: %v\n", err)

	}

	fmt.Printf("Файл успешно сохранён: %s\n", filePath)
	list, _ := FindRowByGroup("changesFile.xlsx", name)
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
