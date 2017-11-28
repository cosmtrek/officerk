package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// Job contains at least one task
type Job struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Schedule  string
	RoutePath string
	Tasks     []Task
}
