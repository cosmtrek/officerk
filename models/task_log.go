package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// TaskStatus has three statuses
type TaskStatus int

// TaskStatus ...
const (
	TaskRunning TaskStatus = iota
	TaskSucceed
	TaskFailed
)

// TaskLog traces every run of a task
type TaskLog struct {
	Comm
	JobLogID uint       `gorm:"not null" json:"joblog_id"`
	TaskID   uint       `gorm:"not null" json:"task_id"`
	Status   TaskStatus `gorm:"not null" json:"status"`
	Retry    int        `json:"-"`
	Result   string     `sql:"type:text" json:"result"`
}
