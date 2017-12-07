package node

import (
	"strings"
	"sync"

	"github.com/cosmtrek/supergo/dag"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type daemon struct {
	db                *gorm.DB
	cron              *cron.Cron
	restartCronCh     chan bool
	restartCronDoneCh chan bool

	sync.RWMutex
	jobs []jobDAG
}

func newDaemon(db *gorm.DB) (*daemon, error) {
	return &daemon{
		db:                db,
		cron:              cron.New(),
		restartCronCh:     make(chan bool),
		restartCronDoneCh: make(chan bool),
	}, nil
}

func (d *daemon) run() error {
	go func() {
		for {
			switch {
			case <-d.restartCronCh:
				d.cron.Stop()
				d.restartCronDoneCh <- true
			}
		}
	}()
	for {
		var err error
		if err = d.reloadCron(); err != nil {
			logrus.Errorf("failed to restart cron, err: %s", err)
		}
		d.cron.Start()
		<-d.restartCronDoneCh
		logrus.Debug("stop running, restarting...")
	}
	return nil
}

func (d *daemon) reloadCron() error {
	var err error
	if err = d.loadJobsFromDB(); err != nil {
		return errors.WithStack(err)
	}
	d.RLock()
	for _, job := range d.jobs {
		logrus.Debugf("job: %s", job.job.Name)
		err = d.cron.AddJob(job.job.Schedule, job)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	d.RUnlock()
	return nil
}

func (d *daemon) restartCron() {
	d.restartCronCh <- true
}

func (d *daemon) loadJobsFromDB() error {
	var err error
	var jobs []Job
	err = d.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&jobs).Error
	if err != nil {
		return err
	}
	var jobDAGs []jobDAG
	for _, job := range jobs {
		var j jobDAG
		if j, err = d.newJobDAG(job); err != nil {
			return err
		}
		jobDAGs = append(jobDAGs, j)
	}
	d.Lock()
	d.jobs = jobDAGs
	d.Unlock()
	return nil
}

func (d *daemon) getJobs() ([]JobRequest, error) {
	var err error
	var jobs []Job
	var tasks []Task
	var data []JobRequest
	err = d.db.Where("deleted_at IS NULL").Order("created_at DESC").Find(&jobs).Error
	if err != nil {
		return data, err
	}
	var taskIDs []uint
	for _, job := range jobs {
		taskIDs = append(taskIDs, job.ID)
	}
	err = d.db.Where("deleted_at IS NULL and id in (?)", taskIDs).Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return data, err
	}
	for _, job := range jobs {
		jr := JobRequest{
			Name:      job.Name,
			Schedule:  job.Schedule,
			RoutePath: job.RoutePath,
			Tasks:     make([]TaskRequest, 0),
		}
		for _, task := range tasks {
			if task.JobID == job.ID {
				jr.Tasks = append(jr.Tasks, TaskRequest{
					Name:      task.Name,
					Command:   task.Command,
					NextTasks: task.NextTasks,
				})
			}
		}
		data = append(data, jr)
	}
	return data, nil
}

func (d *daemon) getJob(id uint) (JobRequest, error) {
	var err error
	var job Job
	var tasks []Task
	var data JobRequest
	err = d.db.Where("deleted_at IS NULL").Where("id = ?", id).Order("created_at DESC").Find(&job).Related(&tasks).Error
	if err != nil {
		return data, err
	}
	data.Name = job.Name
	data.Schedule = job.Schedule
	data.RoutePath = job.RoutePath
	for _, task := range tasks {
		data.Tasks = append(data.Tasks, TaskRequest{
			Name:      task.Name,
			Command:   task.Command,
			NextTasks: task.NextTasks,
		})
	}
	return data, nil
}

func (d *daemon) updateJob(id uint, data Job) error {
	var err error
	var job Job
	var tasks []Task
	if err = d.db.Where("deleted_at IS NULL").Where("id = ?", id).Find(&job).Related(&tasks).Error; err != nil {
		return err
	}
	tx := d.db.Begin()
	// update or create tasks
	for _, taskToUpdate := range job.Tasks {
		for _, taskInDB := range tasks {
			if taskInDB.ID == taskToUpdate.ID {
				if err = d.db.Model(&taskToUpdate).Updates(taskToUpdate).Error; err != nil {
					tx.Rollback()
					return err
				}
				continue
			}
		}
		// new task
		if err = d.db.Create(&taskToUpdate).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// delete old tasks
	for _, taskInDB := range tasks {
		for _, taskToUpdate := range job.Tasks {
			if taskInDB.ID == taskToUpdate.ID {
				continue
			}
		}
		if err = d.db.Delete(&taskInDB).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// update job
	if err = d.db.Model(&job).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (d *daemon) deleteJob(id uint) error {
	var err error
	var job Job
	var tasks []Task
	if err = d.db.Where("deleted_at IS NULL").Where("id = ?", id).Find(&job).Related(&tasks).Error; err != nil {
		return err
	}
	tx := d.db.Begin()
	if err = d.db.Delete(&job).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, task := range tasks {
		if err = d.db.Delete(&task).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

type jobDAG struct {
	db        *gorm.DB
	job       Job
	tasks     []Task
	topoOrder []Task // task name as element in topological order
}

func (d *daemon) newJobDAG(j Job) (jobDAG, error) {
	var err error
	var tasks []Task
	if err = d.db.Model(j).Related(&tasks).Error; err != nil {
		return jobDAG{}, err
	}
	var topoOrder []Task
	if topoOrder, err = getTopologicalOrder(tasks); err != nil {
		return jobDAG{}, err
	}
	return jobDAG{
		db:        d.db,
		job:       j,
		tasks:     tasks,
		topoOrder: topoOrder,
	}, nil
}

// Run implements cron job interface
func (j jobDAG) Run() {
	var err error
	// TODO: if this job is failed, then abort it until it'll be fixed?
	logrus.Debugf("=== run job: %d - %s", j.job.ID, j.job.Name)
	jobLog, err := createJobLog(j.db, j.job.ID)
	if err != nil {
		logrus.Errorf("failed to create job log, err: %s", err)
		return
	}
	for _, t := range j.topoOrder {
		taskLog, err := createTaskLog(j.db, t.JobID, t.ID)
		if err != nil {
			logrus.Errorf("failed to create task log, err: %s", err)
			return
		}
		// TODO: retry?
		if err = runTask(j.db, t); err != nil {
			logrus.Errorf("failed to run task %s, err: %s", t.Name, err)
			if err = updateTaskLogStatus(j.db, taskLog, TaskFailed); err != nil {
				logrus.Errorf("failed to update task log, err: %s", err)
				return
			}
		}
		if err = updateTaskLogStatus(j.db, taskLog, TaskSucceed); err != nil {
			logrus.Errorf("failed to update task log, err: %s", err)
			return
		}
	}
	if err = updateJobLogStatus(j.db, jobLog, JobSucceed); err != nil {
		logrus.Errorf("failed to update job status, err: %s", err)
	}
}

func (d *daemon) runJob(id uint) error {
	d.RLock()
	defer d.RUnlock()
	for _, j := range d.jobs {
		if id == j.job.ID {
			j.Run()
			return nil
		}
	}
	return errors.Errorf("failed to find job %s", id)
}

func getTopologicalOrder(tasks []Task) ([]Task, error) {
	var err error
	dagChecker, err := dag.New(len(tasks))
	if err != nil {
		return nil, err
	}
	for _, tk := range tasks {
		var ntk []dag.ID
		ts := strings.Split(tk.NextTasks, ",")
		for _, t := range ts {
			s := strings.TrimSpace(t)
			if len(s) == 0 {
				continue
			}
			ntk = append(ntk, dag.ID(s))
		}
		v, err := dag.NewVertex(dag.ID(tk.Name), ntk)
		if err != nil {
			return nil, err
		}
		if err = dagChecker.AddVertex(v); err != nil {
			return nil, err
		}
	}
	order := dagChecker.TopologicalOrder()
	if order == nil {
		return nil, errors.New("failed to generate topological order")
	}
	var topoTasks []Task
	for _, e := range order {
		for _, t := range tasks {
			if string(e.ID) == t.Name {
				topoTasks = append(topoTasks, t)
				break
			}
		}
	}
	return topoTasks, nil
}
