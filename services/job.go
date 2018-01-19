package services

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/master/property"
	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/utils"
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
	return db.Where("deleted_at IS NULL").Where("id = ?", id).
		Preload("Node").Preload("Tasks").First(j).Error
}

// GetJobs fetch a list of jobs
func GetJobs(db *gorm.DB) ([]*models.Job, error) {
	var err error
	jobs := make([]*models.Job, 0)
	err = db.Where("deleted_at IS NULL").
		Preload("Tasks").Preload("Node").Find(&jobs).Error
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
	err = db.Where("deleted_at IS NULL").Where("node_id = ?", node.ID).
		Preload("Node").Preload("Tasks").Find(&jobs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return jobs, nil
		}
		return nil, err
	}
	return jobs, nil
}

// RunJobOnNode ...
func RunJobOnNode(j *models.Job, endpoint property.NodeService) error {
	hc := utils.NewHTTPClient()
	url := fmt.Sprintf("http://%s/jobs/%d/run", string(endpoint), j.ID)
	logrus.Debugf("GET: %s", url)
	resp, err := hc.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(string(body))
}

// ReloadJobsOnNode ...
func ReloadJobsOnNode(endpoint property.NodeService) error {
	hc := utils.NewHTTPClient()
	url := fmt.Sprintf("http://%s/reload_jobs", string(endpoint))
	logrus.Debugf("GET: %s", url)
	resp, err := hc.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(string(body))
}

// GetJobLogs ...
func GetJobLogs(db *gorm.DB, j *models.Job) ([]*models.JobLog, error) {
	var err error
	logs := make([]*models.JobLog, 0)
	err = db.Where("deleted_at IS NULL").Where("job_id = ?", j.ID).
		Order("updated_at desc").Limit(20).Find(&logs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return logs, nil
		}
		return nil, err
	}
	return logs, nil
}
