package engine

import (
	_ "time"
	"errors"
)

type Task struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Minutes     int    `json:"minutes"`
	// not sure if necessary
	// StartDate   time.Time `json:"startdate"`
	// Ids are globally unique. Not sure if that's a good idea.
	Id int `json:"id"`
	// not sure if this is a good idea...
	// Parent Task `json:"parent"`
	// subtasks
	SubTasks []*Task `json:"subtasks"`
}

type Person struct {
	Name        string `json:"name"`
	GoalMinutes int    `json:"goalminutes"`
	Tasks       []*Task `json:"tasks"`
}

var PersonStore = make(map[string]*Person)

var taskid int

func NewTask(status, name, description string) Task {
	t := Task{
		Status:      status,
		Name:        name,
		Description: description,
		Id:          taskid,
	}
	taskid++
	return t
}

func FindTask(tasks []*Task, id int) (t *Task, err error) {
	stack := make([]*Task, 0)
	stack = append(stack, tasks...)
	for len(stack) > 0 {
		if stack[0].Id == id {
			t = stack[0]
			return
		} else if len(stack[0].SubTasks) != 0 {
			stack = append(stack, stack[0].SubTasks...)
		}
		// pop
		stack = stack[1:]
	}
	err = errors.New("Task not found")
	return
}
		
