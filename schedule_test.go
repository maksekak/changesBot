package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestHardcodedScheduleAllGroups(t *testing.T) {
	expected := []string{"4-ТМС", "4-ТР", "4-СП", "4-МР", "4-ИС1", "4-ИС2", "4-ИС3", "4-ИС4"}
	for _, g := range expected {
		sched, err := GetHardcodedSchedule(strings.ToLower(g))
		if err != nil {
			t.Fatalf("expected group %s to exist: %v", g, err)
		}
		if sched.Group != g {
			t.Fatalf("expected canonical group %s, got %s", g, sched.Group)
		}
		if len(sched.Days) == 0 {
			t.Fatalf("expected days for group %s", g)
		}
	}
}

func TestFormatDayScheduleWeeksAndChanges(t *testing.T) {
	sched, err := GetHardcodedSchedule("4-ИС3")
	if err != nil {
		t.Fatal(err)
	}

	wedBlue := FormatDaySchedule("4-ИС3", "Среда", "18 ноября", WeekBlue, sched.Days["Среда"], ChangesData{})
	if !strings.Contains(wedBlue, "МДК 01.02") || strings.Contains(wedBlue, "МДК 01.01") {
		t.Fatalf("expected only blue-week MДК 01.02 on Wednesday blue week, got:\n%s", wedBlue)
	}

	wedRed := FormatDaySchedule("4-ИС3", "Среда", "18 ноября", WeekRed, sched.Days["Среда"], ChangesData{})
	if !strings.Contains(wedRed, "МДК 01.01") || strings.Contains(wedRed, "МДК 01.02") {
		t.Fatalf("expected only red-week MДК 01.01 on Wednesday red week, got:\n%s", wedRed)
	}

	sched4, err := GetHardcodedSchedule("4-ИС4")
	if err != nil {
		t.Fatal(err)
	}
	changes := ChangesData{
		DateHeader: "17 ноября (вторник)",
		Entries: []ChangeEntry{
			{PairNumber: 2, Content: "МДК 11.01 — Бобков Д.И. — 307 ауд."},
		},
	}
	tue := FormatDaySchedule("4-ИС4", "Вторник", "17 ноября", WeekBlue, sched4.Days["Вторник"], changes)
	if !strings.Contains(tue, "[ЗАМЕНА: МДК 11.01 — Бобков Д.И. — 307 ауд.]") {
		t.Fatalf("expected replacement on pair 2, got:\n%s", tue)
	}

	// Несовпадающая дата изменений не должна накладываться на другой день
	wedStale := FormatDaySchedule("4-ИС4", "Среда", "18 ноября", WeekBlue, sched4.Days["Среда"], changes)
	if strings.Contains(wedStale, "[ЗАМЕНА:") {
		t.Fatalf("stale changes for 17 ноября should not apply to 18 ноября:\n%s", wedStale)
	}
}

func TestParseChangesFromCachedFile(t *testing.T) {
	if _, err := os.Stat(ChangesCacheFile); err != nil {
		t.Skip("changesSchedule.xlsx not present")
	}
	// Освежаем mtime, чтобы тест не лез в сеть
	now := time.Now()
	_ = os.Chtimes(ChangesCacheFile, now, now)

	data, err := ParseChanges("4-ИС4", "")
	if err != nil {
		t.Fatalf("ParseChanges error: %v", err)
	}
	if !strings.Contains(data.DateHeader, "17 ноября") {
		t.Fatalf("unexpected DateHeader: %q", data.DateHeader)
	}
	if len(data.Entries) != 1 || data.Entries[0].PairNumber != 2 {
		t.Fatalf("unexpected Entries: %+v", data.Entries)
	}
}
