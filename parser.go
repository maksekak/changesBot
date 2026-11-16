package main



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

// ParseMainSchedule возвращает расписание для группы из статических данных (хардкод)
func ParseMainSchedule(groupName string) (GroupSchedule, error) {
	return GetHardcodedSchedule(groupName)
}