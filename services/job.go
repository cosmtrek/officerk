package services

import (
	"github.com/jinzhu/gorm"

	"github.com/cosmtrek/officerk/models"
)

// CreateJob creates job
func CreateJob(db *gorm.DB, j *models.Job) error {
	var err error
	tx := db.Begin()
	// TODO: check node id
	// TODO: check slug
	if err = tx.Create(j).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// UpdateJob updates job
func UpdateJob(db *gorm.DB, j *models.Job, data *models.Job) error {
	var err error
	tx := db.Begin()
	// TODO: check node id
	// TODO: check slug
	if err = tx.Model(j).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Delete useless tasks
	for _, t := range j.Tasks {
		del := true
		for _, t2 := range data.Tasks {
			if t.ID == t2.ID {
				del = false
				break
			}
		}
		if del {
			if err = DeleteTask(db, &t); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

// DeleteJob deletes job
func DeleteJob(db *gorm.DB, j *models.Job) error {
	var err error
	tx := db.Begin()
	for _, t := range j.Tasks {
		if err = DeleteTask(db, &t); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = db.Delete(j).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// GetJob gets job
func GetJob(db *gorm.DB, id string, j *models.Job) error {
	return db.Where("deleted_at IS NULL").Where("id = ?", id).Preload("Tasks").First(j).Error
}

// GetJobs fetch a list of jobs
func GetJobs(db *gorm.DB) ([]*models.Job, error) {
	var err error
	jobs := make([]*models.Job, 0)
	err = db.Where("deleted_at IS NULL").Preload("Tasks").Find(&jobs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return jobs, nil
		}
		return nil, err
	}
	return jobs, nil
}

// GetJobsByNodeIP ...
func GetJobsByNodeIP(db *gorm.DB, ip string) ([]models.Job, error) {
	var err error
	node := new(models.Node)
	jobs := make([]models.Job, 0)
	err = db.Where("deleted_at IS NULL").Where("ip = ?", ip).First(node).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return jobs, nil
		}
		return nil, err
	}
	err = db.Where("deleted_at IS NULL").Where("node_id = ?", node.ID).Preload("Tasks").Find(&jobs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return jobs, nil
		}
		return nil, err
	}
	return jobs, nil
}
