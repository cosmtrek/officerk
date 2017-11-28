package node

import (
	"github.com/jinzhu/gorm"
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
	gorm.Model
	JobLogID uint       `gorm:"not null"`
	TaskID   uint       `gorm:"not null"`
	Name     string     `gorm:"not null"`
	Status   TaskStatus `gorm:"not null"`
	Retry    int
}
