package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ChangeEntry представляет одно изменение
type ChangeEntry struct {
	PairNumber int
	Content    string
}

// ChangesData содержит заголовок даты из таблицы изменений и список изменений группы
type ChangesData struct {
	DateHeader string
	Entries    []ChangeEntry
}

// normalizeChangeContent приводит многострочный текст ячейки изменений к аккуратной строке
func normalizeChangeContent(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	lines := strings.Split(raw, "\n")
	var parts []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			parts = append(parts, line)
		}
	}
	return strings.Join(parts, " — ")
}

// ParseChanges загружает (с кэшированием) и парсит таблицу изменений для указанной группы
// forceRefresh — если "1" или "true", кэш игнорируется и файл скачивается заново
func ParseChanges(groupName string, forceRefresh string) (ChangesData, error) {
	if err := ensureChangesSchedule(forceRefresh == "1" || forceRefresh == "true"); err != nil {
		return ChangesData{}, fmt.Errorf("ошибка кэша изменений: %w", err)
	}

	f, err := excelize.OpenFile(ChangesCacheFile)
	if err != nil {
		return ChangesData{}, fmt.Errorf("ошибка открытия файла изменений: %w", err)
	}
	defer f.Close()

	sheetName, err := getFirstSheetName(f)
	if err != nil {
		return ChangesData{}, fmt.Errorf("ошибка получения листа: %w", err)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return ChangesData{}, fmt.Errorf("ошибка чтения строк изменений: %w", err)
	}

	if len(rows) < 2 {
		return ChangesData{}, nil
	}

	// Первая строка — дата изменений, например "17 ноября (вторник)"
	dateHeader := ""
	if len(rows[0]) > 0 {
		dateHeader = normalizeChangeContent(rows[0][0])
		dateHeader = strings.ReplaceAll(dateHeader, " — ", " ")
	}

	// Находим строку заголовка с номерами пар ("Группа", "1", "2", ...)
	headerIdx := 1
	pairCols := make(map[int]int) // индекс колонки -> номер пары
	for idx, row := range rows {
		if len(row) > 0 && strings.EqualFold(strings.TrimSpace(row[0]), "Группа") {
			headerIdx = idx
			for colIdx := 1; colIdx < len(row); colIdx++ {
				if num, convErr := strconv.Atoi(strings.TrimSpace(row[colIdx])); convErr == nil && num > 0 {
					pairCols[colIdx] = num
				}
			}
			break
		}
	}
	if len(pairCols) == 0 {
		for colIdx := 1; colIdx <= 6; colIdx++ {
			pairCols[colIdx] = colIdx
		}
	}

	// Ищем строку нужной группы (без учёта регистра)
	targetGroup := strings.TrimSpace(groupName)
	groupRow := -1
	for i := headerIdx + 1; i < len(rows); i++ {
		if len(rows[i]) > 0 && strings.EqualFold(strings.TrimSpace(rows[i][0]), targetGroup) {
			groupRow = i
			break
		}
	}

	result := ChangesData{
		DateHeader: dateHeader,
	}

	if groupRow == -1 {
		return result, nil // группы нет в таблице изменений
	}

	row := rows[groupRow]
	for colIdx := 1; colIdx < len(row); colIdx++ {
		pairNum, ok := pairCols[colIdx]
		if !ok {
			continue
		}
		cellContent := normalizeChangeContent(row[colIdx])
		if cellContent == "" {
			continue
		}
		result.Entries = append(result.Entries, ChangeEntry{
			PairNumber: pairNum,
			Content:    cellContent,
		})
	}

	return result, nil
}

// IsCancellation проверяет, является ли изменение отменой
func IsCancellation(content string) bool {
	return strings.Contains(strings.ToLower(content), "отмен")
}
