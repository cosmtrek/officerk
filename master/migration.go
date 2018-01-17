package master

import (
	"github.com/cosmtrek/officerk/models"
)

var tables = []interface{}{
	&models.Node{},
	&models.Job{},
	&models.Task{},
	&models.JobLog{},
	&models.TaskLog{},
}

func (ctr *Controller) autoMigrate() error {
	var err error
	if err = ctr.db.AutoMigrate(tables...).Error; err != nil {
		return err
	}
	err = ctr.db.Model(&models.Task{}).AddUniqueIndex("task_name_and_job_id_unique", "name", "job_id").Error
	if err != nil {
		return err
	}
	return nil
}
