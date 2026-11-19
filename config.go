package main

import "time"

// Ссылки на прямые экспорты Google Таблиц (XLSX)
const (
	MainScheduleURL    = "https://docs.google.com/spreadsheets/d/1IUU0k0st5qsU3sFBfl8wazscE7R6_XHSAlsdJCwUa_4/export?format=xlsx&gid=2041886067"
	ChangesScheduleURL = "https://docs.google.com/spreadsheets/d/1OpDaApQYxlwF9tYmzKbIr_zBt9AO-67aRqKKQXFyU7Q/export?format=xlsx&gid=0"
)

// Имена локальных кэш-файлов
const (
	MainCacheFile    = "mainSchedule.xlsx"
	ChangesCacheFile = "changesSchedule.xlsx"
)

// Настройки кэша и рассылки
const (
	CacheTTL     = 30 * time.Minute // кэш таблицы изменений обновляется не чаще раза в 30 минут
	NotifyHour   = 20               // час вечерней рассылки расписания
	NotifyMinute = 0                // минута вечерней рассылки
	NotifyHour2  = 7                // час утренней рассылки расписания
	NotifyMinute2 = 0               // минута утренней рассылки
)

// Имена файлов настроек
const (
	GroupsFile   = "groups.json"
	WeekStateFile = "week_state.json"
)

// Тип недели
const (
	WeekBlue = "blue"
	WeekRed  = "red"
)

// Цвета заливки ячеек (Excel RGB в формате AARRGGBB)
const (
	ColorRed   = "FFFF0000" // красная неделя
	ColorBlue  = "FF0000FF" // синяя неделя
)