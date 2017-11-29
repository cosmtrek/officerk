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
	Status   TaskStatus `gorm:"not null"`
	Retry    int
}

func createTaskLog(db *gorm.DB, jlID uint, tid uint) (*TaskLog, error) {
	taskLog := &TaskLog{
		JobLogID: jlID,
		TaskID: tid,
		Status: TaskRunning,
	}
	return taskLog, db.Create(taskLog).Error
}

func updateTaskLogStatus(db *gorm.DB, tl *TaskLog, status TaskStatus) error {
	var err error
	if status == TaskSucceed {
		err = db.Model(tl).UpdateColumn("status", status).Error
	} else if status == TaskFailed {
		err = db.Model(tl).UpdateColumn("status", status).UpdateColumn("retry", tl.Retry+1).Error
	}
	return err
}