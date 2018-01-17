package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// Task is the smallest unit for running scripts or commands
type Task struct {
	Comm

	JobID     uint   `gorm:"not null" json:"job_id,omitempty"`
	Name      string `gorm:"not null" json:"name"`
	Command   string `gorm:"not null" json:"command"`
	NextTasks string `json:"next_tasks"` // "task1;task2;task3"
}
