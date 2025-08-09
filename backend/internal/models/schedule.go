package models

import (
	"encoding/json"
	"time"
)

type Schedule interface {
	GetType() string
	JSONify() string
}

type ScheduleForDate struct {
	DaysOff      []string `json:"days_off"`
	Slots        []string `json:"slots"`
	Appointments []string `json:"appointments"`
}

func (s *ScheduleForDate) GetType() string {
	return "slots"
}

func (s *ScheduleForDate) JSONify() string {
	b, _ := json.Marshal(s)
	return string(b)
}

type NextTasks struct {
	Tasks []Task `json:"tasks"`
}

type TaskWithDeadline struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Deadline    time.Time `json:"deadline"`
}

type Task struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
