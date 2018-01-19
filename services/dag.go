package services

import (
	"strings"

	"github.com/cosmtrek/supergo/dag"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/models"
)

// JobDAG ...
type JobDAG struct {
	db        *gorm.DB
	job       models.Job
	topoOrder []models.Task // task name as element in topological order
}

// NewJobDAG ...
func NewJobDAG(db *gorm.DB, j models.Job) (d JobDAG, err error) {
	var topoOrder []models.Task
	if topoOrder, err = getTopologicalOrder(j.Tasks); err != nil {
		return
	}
	d = JobDAG{
		db:        db,
		job:       j,
		topoOrder: topoOrder,
	}
	return
}

// TopoOrder ...
func (j JobDAG) TopoOrder() []models.Task {
	return j.topoOrder
}

// Job ...
func (j JobDAG) Job() models.Job {
	return j.job
}

// Run implements cron job interface
func (j JobDAG) Run() {
	var err error
	// TODO: if this job is failed, then abort it until it'll be fixed?
	logrus.Debugf("=> job[%d]: %s", j.job.ID, j.job.Name)
	jobLog, err := CreateJobLog(j.db, j.job.ID)
	if err != nil {
		logrus.Errorf("failed to create job log, err: %s", err)
		return
	}
	for _, t := range j.topoOrder {
		taskLog, err := CreateTaskLog(j.db, jobLog.ID, t.ID)
		if err != nil {
			logrus.Errorf("failed to create task log, err: %s", err)
			return
		}
		// TODO: retry?
		result, err := RunTask(j.db, t)
		if err != nil {
			logrus.Errorf("failed to run task %s, err: %s", t.Name, err)
			if err = UpdateTaskLogStatus(j.db, taskLog, models.TaskFailed, result); err != nil {
				logrus.Errorf("failed to update task log, err: %s", err)
				return
			}
		}
		if err = UpdateTaskLogStatus(j.db, taskLog, models.TaskSucceed, result); err != nil {
			logrus.Errorf("failed to update task log, err: %s", err)
			return
		}
	}
	if err = UpdateJobLogStatus(j.db, jobLog, models.JobSucceed); err != nil {
		logrus.Errorf("failed to update job status, err: %s", err)
	}
}

func getTopologicalOrder(tasks []models.Task) ([]models.Task, error) {
	var err error
	if len(tasks) == 0 {
		return nil, errors.New("zero tasks")
	}
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
	var topoTasks []models.Task
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
