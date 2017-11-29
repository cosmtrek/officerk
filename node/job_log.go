package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// JobStatus has four statuses following
type JobStatus int

// JobStatus ...
const (
	JobRunning   JobStatus = iota
	JobSucceed
	JobFailed
)

// JobLog traces every run of a job
type JobLog struct {
	gorm.Model
	JobID  uint      `gorm:"not null"`
	Status JobStatus `gorm:"not null"`
	Retry  int
}

func createJobLog(db *gorm.DB, id uint) (*JobLog, error) {
	jl := &JobLog{
		JobID: id,
		Status: JobRunning,
	}
	return jl, db.Create(jl).Error
}

func updateJobLogStatus(db *gorm.DB, jl *JobLog,  status JobStatus) error {
	var err error
	if status == JobSucceed {
		err = db.Model(&jl).UpdateColumn("status", status).Error
	} else if status == JobFailed {
		err = db.Model(&jl).UpdateColumn("status", status).UpdateColumn("retry", jl.Retry+1).Error
	}
	return err
}
