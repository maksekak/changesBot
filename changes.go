package main

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ChangeEntry представляет одно изменение
type ChangeEntry struct {
	PairNumber int
	Content    string
}

// ParseChanges парсит таблицу изменений для указанной группы
func ParseChanges(groupName string, dayOfWeek string) ([]ChangeEntry, error) {
	if err := ensureChangesSchedule(); err != nil {
		return nil, fmt.Errorf("ошибка кэша изменений: %w", err)
	}

	f, err := excelize.OpenFile(ChangesCacheFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла изменений: %w", err)
	}
	defer f.Close()

	sheetName, err := getFirstSheetName(f)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения листа: %w", err)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения строк изменений: %w", err)
	}

	if len(rows) < 3 {
		return nil, nil
	}

	// Находим строку группы в столбце A
	groupRow := -1
	for i, row := range rows {
		if len(row) > 0 && strings.TrimSpace(row[0]) == groupName {
			groupRow = i
			break
		}
	}

	if groupRow == -1 {
		return nil, nil // группы нет в таблице изменений
	}

	// Столбец B-G (индексы 1-6) — пары 1-6
	var changes []ChangeEntry
	groupRowData := rows[groupRow]

	for pairIdx := 1; pairIdx <= 6 && pairIdx < len(groupRowData); pairIdx++ {
		cellContent := strings.TrimSpace(groupRowData[pairIdx])
		if cellContent == "" {
			continue
		}

		changes = append(changes, ChangeEntry{
			PairNumber: pairIdx,
			Content:    cellContent,
		})
	}

	return changes, nil
}

// PairNumberFromLesson извлекает номер пары из строки (для сопоставления)
func PairNumberFromLesson(text string) int {
	// Извлекаем первое число из строки
	var num int
	_, err := fmt.Sscanf(text, "%d", &num)
	if err != nil {
		return 0
	}
	return num
}

// IsCancellation проверяет, является ли изменение отменой
func IsCancellation(content string) bool {
	return strings.Contains(strings.ToLower(content), "отмен")
}

// FormatChange форматирует изменение для вывода
func FormatChange(change ChangeEntry) string {
	if IsCancellation(change.Content) {
		return fmt.Sprintf("%d пара: [ОТМЕНА]", change.PairNumber)
	}
	return fmt.Sprintf("%d пара: %s", change.PairNumber, change.Content)
}