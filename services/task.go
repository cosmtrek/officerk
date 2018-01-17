package services

import (
	"os/exec"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/models"
)

// RunTask executes the task
func RunTask(db *gorm.DB, t models.Task) error {
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

// DeleteTask deletes task
func DeleteTask(db *gorm.DB, t *models.Task) error {
	return db.Delete(t).Error
}
