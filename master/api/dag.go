package api

import (
	"fmt"
	"strings"

	"github.com/cosmtrek/supergo/dag"
	"github.com/pkg/errors"

	"github.com/cosmtrek/officerk/models"
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

type Graph struct {
	Nodes []models.Task `json:"nodes"`
	Edges []edge        `json:"edge"`
}

type edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	ID     string `json:"id"`
}

func newGraph(job *models.Job) *Graph {
	nodes := job.Tasks
	if len(nodes) == 0 {
		return nil
	}
	edges := make([]edge, 0)
	for _, node := range nodes {
		if len(node.NextTasks) == 0 {
			continue
		}
		// TODO: hash task name
		tasks := strings.Split(node.NextTasks, ",")
		for _, task := range tasks {
			edges = append(edges, edge{
				Source: node.Name,
				Target: task,
				ID:     fmt.Sprintf("%s-%s", node.Name, task),
			})
		}
	}
	return &Graph{
		Nodes: nodes,
		Edges: edges,
	}
}
