package engine

import (
	"time"
)

type Task struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Minutes     int       `json:"minutes"`
	StartDate   time.Time `json:"startdate"`
	// Ids are globally unique. Not sure if that's a good idea.
	Id          int       `json:"id"`
}

type Person struct {
	Name        string `json:"name"`
	GoalMinutes int    `json:"goalminutes"`
	Tasks       []Task `json:"done"`
}

var PersonStore = make(map[string]*Person)

var taskid int

func NewTask(status, name string) Task {
	t := Task{Status: status, Name: name, Id: taskid}
	taskid++
	return t
}
