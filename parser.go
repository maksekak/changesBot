package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Lesson представляет одну пару
type Lesson struct {
	Number   int    // номер пары (1-6)
	Time     string // время пары
	Subject  string // предмет
	Teacher  string // преподаватель
	Room     string // аудитория
	Week     string // "blue", "red" или "" (обе недели)
}

// DaySchedule расписание на один день
type DaySchedule struct {
	Day     string
	Lessons []Lesson
}

// GroupSchedule расписание одной группы
type GroupSchedule struct {
	Group  string
	Days   map[string][]Lesson // день недели -> пары
}

// parseCellContent разбивает содержимое ячейки на предмет, преподавателя, аудиторию
func parseCellContent(content string) (subject, teacher, room string) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch i {
		case 0:
			subject = line
		case 1:
			teacher = line
		case 2:
			room = line
		}
	}
	return
}

// getCellRefFromRowCol преобразует индекс строки и столбца в ссылку Excel (A1, B2 и т.д.)
func getCellRefFromRowCol(row int, col int) string {
	colLetter := ""
	c := col + 1
	for c > 0 {
		c--
		colLetter = string(rune('A'+c%26)) + colLetter
		c /= 26
	}
	return fmt.Sprintf("%s%d", colLetter, row+1)
}

// ---------- XML структуры для чтения цветов заливки из XLSX ----------

type xlsxColor struct {
	RGB string `xml:"rgb,attr"`
}

type xlsxPatternFill struct {
	FgColor *xlsxColor `xml:"fgColor"`
	BgColor *xlsxColor `xml:"bgColor"`
}

type xlsxFill struct {
	PatternFill *xlsxPatternFill `xml:"patternFill"`
}

type xlsxFills struct {
	Fill []xlsxFill `xml:"fill"`
}

type xlsxCellXf struct {
	FillId int `xml:"fillId,attr"`
}

type xlsxCellXfs struct {
	CellXf []xlsxCellXf `xml:"cellXf"`
}

type xlsxStyles struct {
	CellXfs xlsxCellXfs `xml:"cellXfs"`
	Fills   xlsxFills   `xml:"fills"`
}

// readFillColors читает styles.xml из XLSX и возвращает карту: индексСтиля -> цветЗаливки (AARRGGBB)
func readFillColors(xlsxPath string) (map[int]string, error) {
	r, err := zip.OpenReader(xlsxPath)
	if err != nil {
		return nil, fmt.Errorf("открытие xlsx: %w", err)
	}
	defer r.Close()

	var stylesFile *zip.File
	for _, f := range r.File {
		if f.Name == "xl/styles.xml" {
			stylesFile = f
			break
		}
	}
	if stylesFile == nil {
		return nil, fmt.Errorf("xl/styles.xml не найден в архиве")
	}

	rc, err := stylesFile.Open()
	if err != nil {
		return nil, fmt.Errorf("открытие styles.xml: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("чтение styles.xml: %w", err)
	}

	var styles xlsxStyles
	if err := xml.Unmarshal(data, &styles); err != nil {
		return nil, fmt.Errorf("парсинг styles.xml: %w", err)
	}

	colorByStyle := make(map[int]string)
	for i, xf := range styles.CellXfs.CellXf {
		if xf.FillId < 0 || xf.FillId >= len(styles.Fills.Fill) {
			continue
		}
		fill := styles.Fills.Fill[xf.FillId]
		if fill.PatternFill == nil {
			continue
		}
		color := ""
		if fill.PatternFill.FgColor != nil && fill.PatternFill.FgColor.RGB != "" {
			color = fill.PatternFill.FgColor.RGB
		} else if fill.PatternFill.BgColor != nil && fill.PatternFill.BgColor.RGB != "" {
			color = fill.PatternFill.BgColor.RGB
		}
		colorByStyle[i] = color
	}

	return colorByStyle, nil
}

// ParseMainSchedule парсит основное расписание из XLSX с определением недели по цвету заливки
func ParseMainSchedule(groupName string) (GroupSchedule, error) {
	if err := ensureMainSchedule(); err != nil {
		return GroupSchedule{}, fmt.Errorf("ошибка кэша основного расписания: %w", err)
	}

	// Читаем цвета заливки из styles.xml
	colorByStyle, err := readFillColors(MainCacheFile)
	if err != nil {
		// Не критично — продолжим без определения недели по цвету
		colorByStyle = make(map[int]string)
	}

	f, err := excelize.OpenFile(MainCacheFile)
	if err != nil {
		return GroupSchedule{}, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	sheetName, err := getFirstSheetName(f)
	if err != nil {
		return GroupSchedule{}, fmt.Errorf("ошибка получения листа: %w", err)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return GroupSchedule{}, fmt.Errorf("ошибка чтения строк: %w", err)
	}

	if len(rows) < 2 {
		return GroupSchedule{}, fmt.Errorf("таблица пуста")
	}

	// Находим колонку для указанной группы в первой строке (заголовки групп)
	groupCol := -1
	headerRow := rows[0]
	for i, cell := range headerRow {
		if strings.TrimSpace(cell) == groupName {
			groupCol = i
			break
		}
	}

	if groupCol == -1 {
		return GroupSchedule{}, fmt.Errorf("группа %s не найдена в заголовках таблицы", groupName)
	}

	schedule := GroupSchedule{
		Group: groupName,
		Days:  make(map[string][]Lesson),
	}

	// Обходим строки таблицы
	currentDay := ""
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// Столбец A (индекс 0) — день недели
		if len(row) > 0 && strings.TrimSpace(row[0]) != "" {
			currentDay = strings.TrimSpace(row[0])
		}

		if currentDay == "" {
			continue
		}

		// Столбец B (индекс 1) — время, C (индекс 2) — номер пары
		if len(row) <= 2 {
			continue
		}

		timeStr := strings.TrimSpace(row[1])
		numberStr := strings.TrimSpace(row[2])

		var pairNumber int
		fmt.Sscanf(numberStr, "%d", &pairNumber)

		// Читаем ячейку для группы
		if len(row) <= groupCol {
			continue
		}
		cellContent := strings.TrimSpace(row[groupCol])
		if cellContent == "" {
			continue
		}

		// Определяем неделя по цвету заливки ячейки
		cellRef := getCellRefFromRowCol(i, groupCol)
		styleID, err := f.GetCellStyle(sheetName, cellRef)
		cellWeek := ""
		if err == nil && styleID >= 0 {
			if color, ok := colorByStyle[styleID]; ok {
				switch color {
				case ColorRed:
					cellWeek = WeekRed
				case ColorBlue:
					cellWeek = WeekBlue
				}
			}
		}

		// Если в ячейке только "ПРАКТИКА" — добавляем как есть
		if cellContent == "ПРАКТИКА" {
			schedule.Days[currentDay] = append(schedule.Days[currentDay], Lesson{
				Number:  pairNumber,
				Time:    timeStr,
				Subject: cellContent,
				Week:    "", // практика — обе недели
			})
			continue
		}

		subject, teacher, room := parseCellContent(cellContent)

		schedule.Days[currentDay] = append(schedule.Days[currentDay], Lesson{
			Number:  pairNumber,
			Time:    timeStr,
			Subject: subject,
			Teacher: teacher,
			Room:    room,
			Week:    cellWeek,
		})
	}

	return schedule, nil
}