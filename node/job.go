package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

var _ xMigrate = new(Job)

// Job contains at least one task
type Job struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Schedule  string
	RoutePath string
}

func (j *Job) migrate(db *gorm.DB) error {
	return nil
}
