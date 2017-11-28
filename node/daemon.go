package node

import (
	"strings"
	"sync"

	"github.com/cosmtrek/supergo/dag"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type daemon struct {
	db *gorm.DB

	sync.RWMutex
	jobs []jobDAG
}

func newDaemon(db *gorm.DB) (*daemon, error) {
	return &daemon{
		db: db,
	}, nil
}

func (d *daemon) run() error {
	return d.loadJobsFromDB()
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

func (d *daemon) runJob(id uint) error {
	d.RLock()
	defer d.RUnlock()
	for _, j := range d.jobs {
		if id == j.job.ID {
			return j.run()
		}
	}
	return errors.Errorf("failed to find job %s", id)
}

func (j *jobDAG) run() error {
	for _, t := range j.topoOrder {
		logrus.Debugf("%d: %s, command: %s", t.JobID, t.Name, t.Command)
	}
	return nil
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
