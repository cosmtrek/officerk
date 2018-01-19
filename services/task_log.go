package services

import (
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/models"
)

// CreateTaskLog creates task log
func CreateTaskLog(db *gorm.DB, jlID uint, tid uint) (*models.TaskLog, error) {
	taskLog := &models.TaskLog{
		JobLogID: jlID,
		TaskID:   tid,
		Status:   models.TaskRunning,
	}
	return taskLog, db.Create(taskLog).Error
}

// UpdateTaskLogStatus updates status of task log
func UpdateTaskLogStatus(db *gorm.DB, tl *models.TaskLog, status models.TaskStatus, result []byte) error {
	var err error
	if status == models.TaskSucceed {
		err = db.Model(tl).UpdateColumn("status", status).
			UpdateColumn("result", string(result)).Error
	} else if status == models.TaskFailed {
		err = db.Model(tl).UpdateColumn("status", status).
			UpdateColumn("retry", tl.Retry+1).
			UpdateColumn("result", string(result)).Error
	}
	return err
}
