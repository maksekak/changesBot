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

// convertToGoogleExportURL преобразует обычный URL Google Таблиц в URL для экспорта в формате XLSX
func convertToGoogleExportURL(originalURL string) (string, error) {
	// Проверяем, что URL содержит нужный домен
	if !strings.Contains(originalURL, "docs.google.com/spreadsheets/d/") {
		return "", fmt.Errorf("URL не является ссылкой на Google Таблицы: %s", originalURL)
	}

	// Извлекаем ID таблицы (между /d/ и следующим / или /edit)
	idStart := strings.Index(originalURL, "/d/") + 3
	if idStart == -1 {
		return "", fmt.Errorf("не удалось извлечь ID таблицы из URL")
	}

	// Ищем позицию /edit или следующего / после /d/
	editPos := strings.Index(originalURL[idStart:], "/edit")
	slashPos := strings.Index(originalURL[idStart:], "/")

	var idEnd int
	if editPos != -1 && slashPos != -1 {
		// Выбираем минимальную позицию между /edit и /
		idEnd = idStart + min(editPos, slashPos)
	} else if editPos != -1 {
		idEnd = idStart + editPos
	} else if slashPos != -1 {
		idEnd = idStart + slashPos
	} else {
		idEnd = len(originalURL)
	}

	tableID := originalURL[idStart:idEnd]

	// Ищем gid (номер листа) после #gid=
	gid := ""
	gidPos := strings.Index(originalURL, "#gid=")
	if gidPos != -1 {
		gid = originalURL[gidPos+5:] // пропускаем "#gid="
		// обрезаем по первому & или концу строки
		andPos := strings.Index(gid, "&")
		if andPos != -1 {
			gid = gid[:andPos]
		}
	}

	// Формируем URL экспорта
	exportURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=xlsx", tableID)
	if gid != "" {
		exportURL += "&gid=" + gid
	}

	return exportURL, nil
}

// Вспомогательная функция для нахождения минимума двух чисел
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func convertToGoogleExportURLForThirdSheet(originalURL string) (string, error) {
	// Проверяем, что URL содержит нужный домен
	if !strings.Contains(originalURL, "docs.google.com/spreadsheets/d/") {
		return "", fmt.Errorf("URL не является ссылкой на Google Таблицы: %s", originalURL)
	}

	// Извлекаем ID таблицы (между /d/ и следующим / или /edit)
	idStart := strings.Index(originalURL, "/d/") + 3
	if idStart == -1 {
		return "", fmt.Errorf("не удалось извлечь ID таблицы из URL")
	}

	// Ищем позицию /edit или следующего / после /d/
	editPos := strings.Index(originalURL[idStart:], "/edit")
	slashPos := strings.Index(originalURL[idStart:], "/")

	var idEnd int
	if editPos != -1 && slashPos != -1 {
		// Выбираем минимальную позицию между /edit и /
		idEnd = idStart + min(editPos, slashPos)
	} else if editPos != -1 {
		idEnd = idStart + editPos
	} else if slashPos != -1 {
		idEnd = idStart + slashPos
	} else {
		idEnd = len(originalURL)
	}

	tableID := originalURL[idStart:idEnd]

	exportURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=xlsx&gid=254196325", tableID)

	return exportURL, nil
}

// DownloadFile скачивает файл и сохраняет его с расширением .xlsx
func DownloadFile(fileURL, outputDir, fileName string) (string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неудачный статус при загрузке: %d", resp.StatusCode)
	}

	// Определяем имя файла из URL (после /export?)

	// Формируем полный путь
	filePath := filepath.Join(outputDir, fileName)

	// Создаём файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Копируем данные
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка записи файла: %w", err)
	}

	return filePath, nil
}
