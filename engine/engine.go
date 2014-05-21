package engine

import (
	"time"
)

type Task struct {
	name        string
	description string
	minutes     int
	date        time.Time
}

type Person struct {
	name        string
	goalMinutes int
	done        []task
}
