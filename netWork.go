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

func FindSecondLinkByText(url, searchText string) (string, error) {
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

	var links []string // будем собирать все подходящие ссылки

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			text := getText(n)
			if strings.Contains(strings.ToLower(text), strings.ToLower(searchText)) {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Проверяем, есть ли хотя бы две ссылки
	if len(links) < 2 {
		return "", fmt.Errorf("найдено меньше двух ссылок с текстом %q (найдено: %d)", searchText, len(links))
	}

	return links[1], nil // возвращаем ВТОРУЮ ссылку (индекс 1)
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

	// Извлекаем ID таблицы (между /d/ и следующим /)
	idStart := strings.Index(originalURL, "/d/") + 3
	if idStart == -1 {
		return "", fmt.Errorf("не удалось извлечь ID таблицы из URL")
	}
	idEnd := strings.Index(originalURL[idStart:], "/")
	if idEnd == -1 {
		idEnd = len(originalURL) - idStart
	} else {
		idEnd += idStart
	}
	tableID := originalURL[idStart:idEnd]

	// Ищем gid (номер листа), если есть
	gid := ""
	gidPos := strings.Index(originalURL, "gid=")
	if gidPos != -1 {
		gid = originalURL[gidPos+4:] // пропускаем "gid="
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

// DownloadFile скачивает файл и сохраняет его с расширением .xlsx
func DownloadFile(fileURL, outputDir string) (string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неудачный статус при загрузке: %d", resp.StatusCode)
	}

	// Определяем имя файла из URL (после /export?)
	fileName := "changesFile.xlsx"

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
