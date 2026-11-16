package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

// getFirstSheetName возвращает имя первого листа в XLSX файле
func getFirstSheetName(f *excelize.File) (string, error) {
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("в файле нет листов")
	}
	return sheets[0], nil
}

// downloadFile скачивает файл по URL и сохраняет локально
func downloadFile(fileURL, outputDir, fileName string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус при загрузке: %d", resp.StatusCode)
	}

	filePath := filepath.Join(outputDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// isCacheFresh проверяет, не устарел ли кэш-файл (младше TTL)
func isCacheFresh(fileName string) bool {
	info, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < CacheTTL
}

// ensureMainSchedule гарантирует, что есть свежая копия основного расписания
func ensureMainSchedule() error {
	if isCacheFresh(MainCacheFile) {
		return nil
	}
	fmt.Println("Кэш основного расписания устарел, скачиваю...")
	return downloadFile(MainScheduleURL, "", MainCacheFile)
}

// ensureChangesSchedule гарантирует, что есть свежая копия таблицы изменений
func ensureChangesSchedule() error {
	if isCacheFresh(ChangesCacheFile) {
		return nil
	}
	fmt.Println("Кэш таблицы изменений устарел, скачиваю...")
	return downloadFile(ChangesScheduleURL, "", ChangesCacheFile)
}