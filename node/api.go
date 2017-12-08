package node

import (
	"strings"

	"github.com/cosmtrek/supergo/dag"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// JobRequest ...
type JobRequest struct {
	Name      string        `json:"name"`
	Schedule  string        `json:"schedule,omitempty"`
	RoutePath string        `json:"route_path,omitempty"`
	Tasks     []TaskRequest `json:"tasks"`
}

// TaskRequest ...
type TaskRequest struct {
	ID        string `json:"id,omitempty"`
	JobID     string `json:"job_id,omitempty"`
	Name      string `json:"name"`
	Command   string `json:"command"`
	NextTasks string `json:"next_tasks,omitempty"` // "task1,task2,task3"
}

func (j *JobRequest) validate() error {
	if j.Schedule == "" && j.RoutePath == "" {
		return errors.New("schedule or route_path must be set")
	}
	if len(j.Tasks) == 0 {
		return errors.New("job must has at least one task")
	}
	return j.checkTasksDependencyCircle()
}

func (j *JobRequest) checkTasksDependencyCircle() error {
	tasks := j.Tasks
	taskVertices := make(map[string][]string, len(tasks))
	taskCache := make(map[string]bool, len(tasks))

	for _, tk := range tasks {
		var ntk []string
		ts := strings.Split(tk.NextTasks, ",")
		for _, t := range ts {
			s := strings.TrimSpace(t)
			if len(s) == 0 {
				continue
			}
			ntk = append(ntk, s)
		}
		taskVertices[tk.Name] = ntk
		taskCache[tk.Name] = true
	}
	for tk, ntk := range taskVertices {
		for _, s := range ntk {
			if _, ok := taskCache[s]; !ok {
				return errors.Errorf("failed to find %s in %s", s, tk)
			}
		}
	}

	var err error
	dagChecker, err := dag.New(len(tasks))
	if err != nil {
		return err
	}
	for tk, ntk := range taskVertices {
		var nID []dag.ID
		for _, s := range ntk {
			nID = append(nID, dag.ID(s))
		}
		vertex, err := dag.NewVertex(dag.ID(tk), nID)
		if err != nil {
			return err
		}
		if err = dagChecker.AddVertex(vertex); err != nil {
			return err
		}
	}
	if hasCircle := dagChecker.CheckCircle(); hasCircle {
		return errors.Errorf("found circle in this job, circle path: %s", dagChecker.CirclePath())
	}
	return nil
}

func (j *JobRequest) save(db *gorm.DB) error {
	var err error
	tx := db.Begin()
	job := Job{
		Name:      j.Name,
		Schedule:  j.Schedule,
		RoutePath: j.RoutePath,
	}
	if err = tx.Create(&job).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, task := range j.Tasks {
		task := Task{
			JobID:     job.ID,
			Name:      task.Name,
			Command:   task.Command,
			NextTasks: task.NextTasks,
		}
		if err = tx.Create(&task).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
