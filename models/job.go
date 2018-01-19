package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// JobType has two types
type JobType int

// JobType ...
const (
	JobTypeCron JobType = iota
	JobTypeManual
)

// Job contains at least one task
type Job struct {
	Comm
	Name     string `gorm:"not null" json:"name"`
	Typ      string `gorm:"not null" json:"typ"`
	Schedule string `json:"schedule,omitempty"`
	Slug     string `json:"slug,omitempty"`
	NodeID   uint   `gorm:"not null" json:"-"`

	Node  Node   `json:"node"`
	Tasks []Task `gorm:"ForeignKey:JobID" json:"tasks"`
}
