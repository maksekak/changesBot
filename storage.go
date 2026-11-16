package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GroupSettings хранит настройки для одного чата
type GroupSettings struct {
	GroupName string `json:"group_name"`
	Notify    bool   `json:"notify"`
}

// WeekState хранит глобальное состояние недели
type WeekState struct {
	ReferenceDate string `json:"reference_date"` // дата отсчёта в формате YYYY-MM-DD
	ReferenceWeek string `json:"reference_week"` // blue или red
}

// loadGroups загружает настройки групп из файла
func loadGroups() map[int64]GroupSettings {
	data, err := os.ReadFile(GroupsFile)
	if err != nil {
		return make(map[int64]GroupSettings)
	}
	var groups map[int64]GroupSettings
	if err := json.Unmarshal(data, &groups); err != nil {
		return make(map[int64]GroupSettings)
	}
	return groups
}

// saveGroups сохраняет настройки групп в файл
func saveGroups(groups map[int64]GroupSettings) error {
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GroupsFile, data, 0644)
}

// setGroup устанавливает группу для чата
func setGroup(chatID int64, groupName string) error {
	groups := loadGroups()
	groups[chatID] = GroupSettings{GroupName: groupName, Notify: groups[chatID].Notify}
	return saveGroups(groups)
}

// getGroup возвращает группу для чата
func getGroup(chatID int64) (string, bool) {
	groups := loadGroups()
	g, ok := groups[chatID]
	if !ok {
		return "", false
	}
	return g.GroupName, true
}

// toggleNotify переключает уведомления для чата
func toggleNotify(chatID int64) bool {
	groups := loadGroups()
	g, ok := groups[chatID]
	if !ok {
		groups[chatID] = GroupSettings{GroupName: "", Notify: true}
		saveGroups(groups)
		return true
	}
	g.Notify = !g.Notify
	groups[chatID] = g
	saveGroups(groups)
	return g.Notify
}

// isNotify возвращает true, если уведомления включены для чата
func isNotify(chatID int64) bool {
	groups := loadGroups()
	g, ok := groups[chatID]
	return ok && g.Notify
}

// loadWeekState загружает состояние недели
func loadWeekState() WeekState {
	data, err := os.ReadFile(WeekStateFile)
	if err != nil {
		return WeekState{}
	}
	var ws WeekState
	if err := json.Unmarshal(data, &ws); err != nil {
		return WeekState{}
	}
	return ws
}

// saveWeekState сохраняет состояние недели
func saveWeekState(ws WeekState) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(WeekStateFile, data, 0644)
}

// setWeek устанавливает текущую неделю с привязкой к сегодняшней дате (обратная совместимость)
func setWeek(week string) error {
	return setWeekWithDate(time.Now().Format("2006-01-02"), week)
}

// setWeekWithDate устанавливает неделю с заданной датой отсчёта (формат YYYY-MM-DD)
func setWeekWithDate(dateStr string, week string) error {
	ws := WeekState{
		ReferenceDate: dateStr,
		ReferenceWeek: week,
	}
	return saveWeekState(ws)
}

// parseDateDDMM парсит дату в формате ДД.ММ и возвращает YYYY-MM-DD (текущий год)
func parseDateDDMM(s string) (string, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("ожидается формат ДД.ММ, например 01.09")
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return "", fmt.Errorf("некорректный день: %s", parts[0])
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return "", fmt.Errorf("некорректный месяц: %s", parts[1])
	}
	year := time.Now().Year()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day), nil
}

// getMonday возвращает дату понедельника недели, к которой относится переданная дата
func getMonday(t time.Time) time.Time {
	weekday := int(t.Weekday()) // 0 = воскресенье, 1 = понедельник, ..., 6 = суббота
	// Сдвигаем назад до понедельника (если воскресенье, отнимаем 6 дней)
	daysToSub := weekday - 1
	if daysToSub < 0 {
		daysToSub = 6
	}
	return t.AddDate(0, 0, -daysToSub)
}

// getCurrentWeek вычисляет текущую неделю с привязкой к понедельникам
func getCurrentWeek() (string, error) {
	ws := loadWeekState()
	if ws.ReferenceDate == "" || ws.ReferenceWeek == "" {
		return "", fmt.Errorf("неделя не установлена. Используйте /weekset ДД.ММ blue|red, например /weekset 01.09 blue")
	}

	refDate, err := time.Parse("2006-01-02", ws.ReferenceDate)
	if err != nil {
		return "", err
	}

	now := time.Now()

	refMonday := getMonday(refDate)
	nowMonday := getMonday(now)

	daysBetween := int(nowMonday.Sub(refMonday).Hours() / 24)
	weeksBetween := daysBetween / 7

	if weeksBetween%2 == 0 {
		return ws.ReferenceWeek, nil
	}

	if ws.ReferenceWeek == WeekBlue {
		return WeekRed, nil
	}
	return WeekBlue, nil
}