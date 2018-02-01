package services

import (
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/models"
)

// CreateJobLog creates job log
func CreateJobLog(db *gorm.DB, id uint) (*models.JobLog, error) {
	jl := &models.JobLog{
		JobID:  id,
		Status: models.JobRunning,
	}
	return jl, db.Create(jl).Error
}

// UpdateJobLogStatus updates status of job log
func UpdateJobLogStatus(db *gorm.DB, jl *models.JobLog, status models.JobStatus) error {
	var err error
	if status == models.JobSucceed {
		err = db.Model(&jl).UpdateColumn("status", status).Error
	} else if status == models.JobFailed {
		err = db.Model(&jl).UpdateColumn("status", status).UpdateColumn("retry", jl.Retry+1).Error
	}
	return err
}

// GetJobLog gets job log
func GetJobLog(db *gorm.DB, id string, l *models.JobLog) error {
	return db.Where("deleted_at IS NULL").Where("id = ?", id).
		Preload("Job").Preload("TaskLogs").Preload("TaskLogs.Task").First(l).Error
}

// ListJobLogs gets latest job log
func ListJobLogs(db *gorm.DB) ([]*models.JobLog, error) {
	var err error
	logs := make([]*models.JobLog, 0)
	err = db.Where("deleted_at IS NULL").Order("updated_at desc").Limit(20).
		Preload("Job").Preload("TaskLogs").Find(&logs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return logs, nil
		}
		return nil, err
	}
	return logs, nil
}
