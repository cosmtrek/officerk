package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

var _ xMigrate = new(TaskLog)

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
	JobLogID int        `gorm:"not null"`
	TaskID   int        `gorm:"not null"`
	Name     string     `gorm:"not null"`
	Status   TaskStatus `gorm:"not null"`
	Retry    int
}

func (tl *TaskLog) migrate(db *gorm.DB) error {
	return nil
}
