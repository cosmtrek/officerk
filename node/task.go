package node

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
)

// Task is the smallest unit for running scripts or commands
type Task struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Command   string `gorm:"not null"`
	NextTasks string // "task1;task2;task3"
}
