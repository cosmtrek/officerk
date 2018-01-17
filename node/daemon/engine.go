package daemon

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/services"
	"github.com/cosmtrek/officerk/utils"
)

// Engine ...
type Engine struct {
	ip                string // where node runs
	db                *gorm.DB
	cron              *cron.Cron
	restartCronCh     chan bool
	restartCronDoneCh chan bool

	sync.RWMutex
	jobDAGs []services.JobDAG
}

// NewEngine ...
func NewEngine(db *gorm.DB) (*Engine, error) {
	ip, err := utils.GetIP()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("Node IP: %s", ip)
	return &Engine{
		ip:                ip,
		db:                db,
		cron:              cron.New(),
		restartCronCh:     make(chan bool),
		restartCronDoneCh: make(chan bool),
	}, nil
}

// Run ...
func (e *Engine) Run() error {
	go func() {
		for {
			switch {
			case <-e.restartCronCh:
				e.cron.Stop()
				e.restartCronDoneCh <- true
			}
		}
	}()
	for {
		var err error
		if err = e.reloadCron(); err != nil {
			logrus.Errorf("failed to restart cron, err: %s", err)
		}
		e.cron.Start()
		<-e.restartCronDoneCh
		logrus.Debug("stop running cron, restarting...")
	}
	return nil
}

// Reload restarts daemon
func (e *Engine) Reload() {
	e.restartCronCh <- true
}

// RunJob ...
func (e *Engine) RunJob(id uint) error {
	e.RLock()
	defer e.RUnlock()
	for _, j := range e.jobDAGs {
		if id == j.Job().ID {
			j.Run()
			return nil
		}
	}
	return errors.Errorf("failed to find job %s", id)
}

func (e *Engine) reloadCron() error {
	var err error
	if err = e.fetchJobs(); err != nil {
		return errors.WithStack(err)
	}
	e.Lock()
	ncron := cron.New()
	for _, dag := range e.jobDAGs {
		logrus.Debugf("job: %s", dag.Job().Name)
		err = ncron.AddJob(dag.Job().Schedule, dag)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	e.cron.Stop()
	e.cron = ncron
	e.Unlock()
	return nil
}

func (e *Engine) fetchJobs() error {
	var err error
	jobs, err := services.GetJobsByNodeIP(e.db, e.ip)
	if err != nil {
		return err
	}
	dags := make([]services.JobDAG, 0)
	for _, job := range jobs {
		d, err := services.NewJobDAG(e.db, job)
		if err != nil {
			return err
		}
		dags = append(dags, d)
	}
	e.Lock()
	e.jobDAGs = dags
	e.Unlock()
	return nil
}
