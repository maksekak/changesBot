package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func FindDoc(url, searchText string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP ошибка: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	var link string // будем хранить первую найденную подходящую ссылку

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			text := getText(n)
			if strings.Contains(strings.ToLower(text), strings.ToLower(searchText)) {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link = attr.Val // сохраняем первую найденную ссылку
						return          // прерываем поиск после нахождения первой ссылки
					}
				}
			}
		}
		// продолжаем обход дочерних узлов, только если ссылка ещё не найдена
		if link == "" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)

	// Проверяем, найдена ли ссылка
	if link == "" {
		return "", fmt.Errorf("не найдено ссылок с текстом %q", searchText)
	}

	return link, nil // возвращаем найденную ссылку
}

// getText собирает текст из узла
func getText(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text += n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return strings.TrimSpace(text)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func convertToGoogleExportURL(originalURL string, sheetID int) (string, error) {
	// Проверяем, что URL содержит нужный домен
	if !strings.Contains(originalURL, "docs.google.com/spreadsheets/d/") {
		return "", fmt.Errorf("URL не является ссылкой на Google Таблицы: %s", originalURL)
	}

	// Находим позицию начала ID таблицы (/d/)
	idStart := strings.Index(originalURL, "/d/") + len("/d/")
	if idStart == -1 {
		return "", fmt.Errorf("не удалось извлечь ID таблицы из URL")
	}

	// Ищем позицию /edit или следующего / после /d/
	remainingURL := originalURL[idStart:]
	editPos := strings.Index(remainingURL, "/edit")
	slashPos := strings.Index(remainingURL, "/")

	var idEnd int
	switch {
	case editPos != -1 && slashPos != -1:
		// Выбираем минимальную позицию между /edit и /
		idEnd = idStart + min(editPos, slashPos)
	case editPos != -1:
		idEnd = idStart + editPos
	case slashPos != -1:
		idEnd = idStart + slashPos
	default:
		idEnd = len(originalURL)
	}

	tableID := originalURL[idStart:idEnd]

	// Формируем URL экспорта с указанным sheetID (gid)
	exportURL := fmt.Sprintf(
		"https://docs.google.com/spreadsheets/d/%s/export?format=xlsx&gid=%d",
		tableID,
		sheetID,
	)

	return exportURL, nil
}

// DownloadFile скачивает файл и сохраняет его с расширением .xlsx
func DownloadFile(fileURL, outputDir, fileName string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус при загрузке: %d", resp.StatusCode)
	}

	// Определяем имя файла из URL (после /export?)

	// Формируем полный путь
	filePath := filepath.Join(outputDir, fileName)

	// Создаём файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Копируем данные
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}
func GetSchedule(url, linkName, fileName string, sID int) {
	excelLink, err := FindDoc(url, linkName)
	if err != nil {
		fmt.Printf("Ошибка поиска второй ссылки: %v\n", err)

	}
	fmt.Printf("Найденная вторая ссылка: %s\n", excelLink)

	// 2. Преобразуем URL Google Таблиц в URL экспорта
	exportURL, err := convertToGoogleExportURL(excelLink, sID)
	if err != nil {
		fmt.Printf("Ошибка преобразования URL в формат экспорта: %v\n", err)

	}
	fmt.Printf("URL для экспорта: %s\n", exportURL)

	// 3. Скачиваем файл
	err = DownloadFile(exportURL, "", fileName)
	if err != nil {
		fmt.Printf("Ошибка скачивания файла: %v\n", err)

	}
}
