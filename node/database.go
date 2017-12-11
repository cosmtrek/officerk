package node

import (
	"github.com/jinzhu/gorm"
)

func dbFetchJobAndTasks(db *gorm.DB, id uint, job *Job, tasks *[]Task) error {
	var err error
	err = db.Where("deleted_at IS NULL").Where("id = ?", id).First(job).Error
	if err != nil {
		return err
	}
	err = db.Where("deleted_at IS NULL").Where("job_id = ?", job.ID).Order("created_at ASC").Find(tasks).Error
	if err != nil {
		return err
	}
	return nil
}
