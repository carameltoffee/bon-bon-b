package dates

import (
	"fmt"
	"strings"
)

var weekdayMap = map[string]int{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

func WeekdayToInt(day string) (int, error) {
	day = strings.ToLower(day)
	if val, ok := weekdayMap[day]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("invalid weekday: %s", day)
}

func IntToWeekday(i int) string {
	for k, v := range weekdayMap {
		if v == i {
			return k
		}
	}
	return "unknown"
}
