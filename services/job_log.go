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