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

type jobDAG struct {
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
		job:       j,
		tasks:     tasks,
		topoOrder: topoOrder,
	}, nil
}

// Run implements cron job interface
func (j jobDAG) Run() {
	for _, t := range j.topoOrder {
		logrus.Debugf("%d: %s, command: %s", t.JobID, t.Name, t.Command)
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
