package api

import (
	"strings"

	"github.com/cosmtrek/supergo/dag"
	"github.com/pkg/errors"
)

func checkTasksDependencyCircle(j *JobRequest) error {
	tasks := j.Job.Tasks
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
