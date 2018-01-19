package services

import (
	"os/exec"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/models"
)

// RunTask executes the task
func RunTask(db *gorm.DB, t models.Task) ([]byte, error) {
	var c *exec.Cmd
	c = exec.Command("/bin/sh", "-c", t.Command)
	logrus.Debugf("-> task: %s", t.Name)
	return c.CombinedOutput()
}

// DeleteTask deletes task
func DeleteTask(db *gorm.DB, t *models.Task) error {
	return db.Delete(t).Error
}
