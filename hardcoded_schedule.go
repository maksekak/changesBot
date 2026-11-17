package main

import (
	"fmt"
	"sort"
	"strings"
)

// HardcodedSchedule содержит точное расписание всех групп из основной таблицы (лист "4 Курс").
// Поле Week: "" — обе недели, "blue" — только синяя, "red" — только красная.
var HardcodedSchedule = map[string]map[string][]Lesson{
	"4-ТМС": {
		"Понедельник": {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Вторник":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Среда":       {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Четверг":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Пятница":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
	},
	"4-ТР": {
		"Понедельник": {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Вторник":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Среда":       {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Четверг":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Пятница":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
	},
	"4-СП": {
		"Понедельник": {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Вторник":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Среда":       {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Четверг":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Пятница":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
	},
	"4-МР": {
		"Понедельник": {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Вторник":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Среда":       {{Number: 3, Time: "12:45 - 14:20", Subject: "ПРАКТИКА", Week: ""}},
		"Четверг":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
		"Пятница":     {{Number: 3, Time: "12:20 - 13:55", Subject: "ПРАКТИКА", Week: ""}},
	},
	"4-ИС1": {
		"Понедельник": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Разговоры о важном", Teacher: "Бодров Е.Н.", Room: "410", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
		"Вторник": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Менеджмент в ПД", Teacher: "Рубан Е.А.", Room: "106", Week: "red"},
		},
		"Среда": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "Физра", Teacher: "Видьманов А.А.", Room: "с/з", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Семьеведение", Teacher: "", Room: "", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 5, Time: "16:10 - 17:40", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: "blue"},
		},
		"Четверг": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "Менеджмент ", Teacher: "Рубан Е.А.", Room: "101", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: "red"},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: "blue"},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Иностранный язык", Teacher: "Васильева В.В., Староверова Г.Б.", Room: "2, 1", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: "red"},
			{Number: 5, Time: "15:45 - 17:15", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: "blue"},
		},
		"Пятница": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "307", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.03", Teacher: "Авдоян Д.Т.", Room: "306", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
	},
	"4-ИС2": {
		"Понедельник": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Разговоры о важном", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 01.03", Teacher: "Авдоян Д.Т.", Room: "306", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
		"Вторник": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "Менеджмент в ПД", Teacher: "Рубан Е.А.", Room: "", Week: "blue"},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: "red"},
		},
		"Среда": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Семьеведение", Teacher: "", Room: "", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "Физра", Teacher: "Видьманов А.А.", Room: "с/з", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
		"Четверг": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: "red"},
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: "blue"},
			{Number: 3, Time: "12:20 - 13:55", Subject: "Иностранный язык", Teacher: "Анохина С.О., Зуева Н.Н.", Room: "13, 25", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: "red"},
			{Number: 5, Time: "15:45 - 17:15", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: "blue"},
		},
		"Пятница": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "Менеджмент ", Teacher: "Рубан Е.А.", Room: "101", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "307", Week: ""},
		},
	},
	"4-ИС3": {
		"Понедельник": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Разговоры о важном", Teacher: "Шведко Я.В.", Room: "307", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "Менеджмент ", Teacher: "Рубан Е.А.", Room: "101", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "Иностранный язык", Teacher: "Анохина С.О., Тормосина Т.М.", Room: "13, 36-о", Week: ""},
			{Number: 5, Time: "16:10 - 17:40", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: "blue"},
		},
		"Вторник": {
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.03", Teacher: "Авдоян Д.Т.", Room: "306", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
		"Среда": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Семьеведение", Teacher: "", Room: "", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "Физра", Teacher: "Видьманов А.А.", Room: "с/з", Week: ""},
			{Number: 5, Time: "16:10 - 17:40", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: "red"},
			{Number: 5, Time: "16:10 - 17:40", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: "blue"},
		},
		"Четверг": {
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "106", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: ""},
		},
		"Пятница": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "Менеджмент в ПД", Teacher: "Рубан Е.А.", Room: "101", Week: "red"},
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "307", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: "red"},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: "blue"},
		},
	},
	"4-ИС4": {
		"Понедельник": {
			{Number: 0, Time: "12:20 - 12:35", Subject: "Разговоры о важном", Teacher: "Петрушенко Л.Л.", Room: "310", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: ""},
			{Number: 5, Time: "16:10 - 17:40", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
		"Вторник": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "Физра", Teacher: "Видьманов А.А.", Room: "с/з", Week: ""},
			{Number: 3, Time: "12:20 - 13:55", Subject: "Менеджмент в ПД", Teacher: "Рубан Е.А.", Room: "106", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: ""},
		},
		"Среда": {
			{Number: 1, Time: "8:30 - 10:05", Subject: "Иностранный язык", Teacher: "Зуева Н.Н., Тормосина Т.М.", Room: "25, 36-о", Week: ""},
			{Number: 2, Time: "10:15 - 11:50", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 0, Time: "12:20 - 12:35", Subject: "Семьеведение", Teacher: "", Room: "", Week: ""},
			{Number: 3, Time: "12:45 - 14:20", Subject: "МДК 01.03", Teacher: "Авдоян Д.Т.", Room: "306", Week: ""},
			{Number: 4, Time: "14:30 - 16:00", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
		},
		"Четверг": {
			{Number: 3, Time: "12:20 - 13:55", Subject: "МДК 01.04", Teacher: "Евстигнеев А.С.", Room: "209", Week: ""},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: ""},
			{Number: 5, Time: "15:45 - 17:15", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
			{Number: 6, Time: "17:25 - 18:55", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: "red"},
			{Number: 6, Time: "17:25 - 18:55", Subject: "МДК 01.02", Teacher: "Ишкова Л.Г.", Room: "305", Week: "blue"},
		},
		"Пятница": {
			{Number: 3, Time: "12:20 - 13:55", Subject: "Экономика отрасли", Teacher: "Кондратюк Е.В.", Room: "405", Week: "red"},
			{Number: 3, Time: "12:20 - 13:55", Subject: "Менеджмент в ПД", Teacher: "Рубан Е.А.", Room: "308", Week: "blue"},
			{Number: 4, Time: "14:05 - 15:35", Subject: "Стандартизация", Teacher: "Петрушенко Л.Л.", Room: "310", Week: "red"},
			{Number: 4, Time: "14:05 - 15:35", Subject: "МДК 11.01", Teacher: "Бобков Д.И.", Room: "302", Week: "blue"},
			{Number: 5, Time: "15:45 - 17:15", Subject: "МДК 01.01", Teacher: "Ишкова Л.Г.", Room: "305", Week: ""},
		},
	},
}

// GetHardcodedSchedule возвращает расписание для группы из статических данных (без учёта регистра)
func GetHardcodedSchedule(groupName string) (GroupSchedule, error) {
	target := strings.TrimSpace(groupName)
	for name, days := range HardcodedSchedule {
		if strings.EqualFold(name, target) {
			schedule := GroupSchedule{
				Group: name,
				Days:  make(map[string][]Lesson),
			}
			for day, lessons := range days {
				schedule.Days[day] = lessons
			}
			return schedule, nil
		}
	}
	return GroupSchedule{}, fmt.Errorf("группа %s не найдена в расписании. Доступные группы: %s", groupName, strings.Join(GetAvailableGroups(), ", "))
}

// GetAvailableGroups возвращает отсортированный список доступных групп
func GetAvailableGroups() []string {
	groups := make([]string, 0, len(HardcodedSchedule))
	for name := range HardcodedSchedule {
		groups = append(groups, name)
	}
	sort.Strings(groups)
	return groups
}
