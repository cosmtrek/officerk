package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// JobStatus has three statuses
type JobStatus int

// JobStatus ...
const (
	JobRunning JobStatus = iota
	JobSucceed
	JobFailed
)

// JobLog traces every run of a job
type JobLog struct {
	Comm
	JobID  uint      `gorm:"not null" json:"job_id"`
	Status JobStatus `gorm:"not null" json:"status"`
	Retry  int
}
