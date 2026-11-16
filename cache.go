package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

var cacheMu sync.Mutex

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// getFirstSheetName возвращает имя первого листа в XLSX файле
func getFirstSheetName(f *excelize.File) (string, error) {
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("в файле нет листов")
	}
	return sheets[0], nil
}

// downloadFile скачивает файл по URL и сохраняет локально атомарно через временный файл
func downloadFile(fileURL, fileName string) error {
	resp, err := httpClient.Get(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус при загрузке: %d", resp.StatusCode)
	}

	tmpFile := fileName + ".tmp"
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		file.Close()
		os.Remove(tmpFile)
		return fmt.Errorf("ошибка записи файла: %w", err)
	}
	file.Close()

	if err := os.Rename(tmpFile, fileName); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("ошибка сохранения файла: %w", err)
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

// ensureChangesSchedule гарантирует, что есть свежая копия таблицы изменений
// (потокобезопасно, чтобы при массовой рассылке или одновременных запросах был только 1 запрос к Google)
func ensureChangesSchedule() error {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if isCacheFresh(ChangesCacheFile) {
		return nil
	}
	fmt.Println("Кэш таблицы изменений устарел, скачиваю...")
	return downloadFile(ChangesScheduleURL, ChangesCacheFile)
}
