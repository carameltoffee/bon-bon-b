package models

import "time"

type SchedulePattern struct {
	UserId        int64         `json:"user_id"`
	Type          string        `json:"schedule_type"`
	Slots         []SlotPattern `json:"time_slots"`
	DaysAhead     int           `json:"days_ahead"`
	LastGenerated time.Time     `json:"last_generated"`
}

type SlotPattern struct {
	Weekday string `json:"weekday"`
	Time    string `json:"time"`
}
