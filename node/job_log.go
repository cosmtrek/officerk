package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

var _ xMigrate = new(JobLog)

// JobStatus has four statuses following
type JobStatus int

// JobStatus ...
const (
	JobInit JobStatus = iota
	JobRunning
	JobSucceed
	JobFailed
)

// JobLog traces every run of a job
type JobLog struct {
	gorm.Model
	JobID  int       `gorm:"not null"`
	Status JobStatus `gorm:"not null"`
	Retry  int
}

func (jl *JobLog) migrate(db *gorm.DB) error {
	return nil
}
