package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// Node ...
type Node struct {
	Comm
	Name string `gorm:"not null" json:"name"`
	IP   string `gorm:"not null" json:"ip"`

	Jobs []Job `json:"jobs,omitempty"`
}
