package node

import (
	"os/exec"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"github.com/sirupsen/logrus"
)

// Task is the smallest unit for running scripts or commands
type Task struct {
	gorm.Model
	JobID     uint   `gorm:"not null"`
	Name      string `gorm:"not null"`
	Command   string `gorm:"not null"`
	NextTasks string // "task1;task2;task3"
}

func runTask(db *gorm.DB, t Task) error {
	var err error
	var c *exec.Cmd
	c = exec.Command("/bin/sh", "-c", t.Command)
	logrus.Debugf("run task: %s", t.Name)
	output, err := c.CombinedOutput()
	if err != nil {
		return err
	}
	logrus.Debugf("%s output: %s", t.Name, string(output))
	return err
}
