package engine

import (
	_ "encoding/json"
	"time"
)

type Task struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Minutes     int       `json:"minutes"`
	StartDate   time.Time `json:"startdate"`
}

type Person struct {
	Name        string `json:"name"`
	GoalMinutes int    `json:"goalminutes"`
	Done        []Task `json:"done"`
}

type PersonStore map[string]Person

func NewTask(status, name string) Task {
	return Task{Status: status, Name: name}
}
