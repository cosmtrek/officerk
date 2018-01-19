package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/cosmtrek/officerk/master/property"
	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/services"
	"github.com/cosmtrek/officerk/utils/api"
)

type ctxKey string

var jobKey = ctxKey("job")

// JobRequest ...
type JobRequest struct {
	*models.Job
}

// Bind for post-processing JobRequest
func (j *JobRequest) Bind(r *http.Request) error {
	var err error
	if len(j.Tasks) == 0 {
		return errors.New("job must has at least one task")
	}
	if err = checkTasksDependencyCircle(j); err != nil {
		return err
	}
	ts := make([]models.Task, 0)
	for _, t := range j.Tasks {
		t2 := t
		t2.CreatedAt = time.Now()
		t2.UpdatedAt = time.Now()
		ts = append(ts, t2)
	}
	j.Tasks = ts
	return nil
}

// JobResponse ...
type JobResponse struct {
	*models.Job
	Graph *Graph `json:"graph,omitempty"`
}

// NewJobResponse ...
func NewJobResponse(job *models.Job) *JobResponse {
	return &JobResponse{Job: job, Graph: newGraph(job)}
}

// Render for JobResponse
func (jr JobResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// JobListResponse ...
type JobListResponse []*JobResponse

// NewJobListResponse ...
func NewJobListResponse(jobs []*models.Job) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, job := range jobs {
		list = append(list, NewJobResponse(job))
	}
	return list
}

// JobCtx finds job
func JobCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id := chi.URLParam(r, "jobID")
		if id == "" {
			render.Render(w, r, api.ErrNotFound)
			return
		}
		job := new(models.Job)
		if err = services.GetJob(db, id, job); err != nil {
			render.Render(w, r, api.ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), jobKey, job)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ListJobs finds all jobs
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	var err error
	jobs, err := services.GetJobs(db)
	if err != nil {
		render.Render(w, r, api.ErrNotFound)
		return
	}
	render.Render(w, r, api.OK(NewJobListResponse(jobs)))
}

// CreateJob creates job
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var err error
	data := &JobRequest{}
	if err = render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.CreateJob(db, data.Job); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = h.reloadJobsOnNode(data.Job); err != nil {
		render.Render(w, r, api.ErrNodeResponse(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, api.OK(NewJobResponse(data.Job)))
}

// GetJob find the job
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	job := r.Context().Value(jobKey).(*models.Job)
	render.Render(w, r, api.OK(NewJobResponse(job)))
}

// UpdateJob updates the job
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	var err error
	job := r.Context().Value(jobKey).(*models.Job)
	data := &JobRequest{}
	if err = render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.UpdateJob(db, job, data.Job); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = h.reloadJobsOnNode(job); err != nil {
		render.Render(w, r, api.ErrNodeResponse(err))
		return
	}
	render.Render(w, r, api.OK(NewJobResponse(data.Job)))
}

// DeleteJob deletes the job
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	var err error
	job := r.Context().Value(jobKey).(*models.Job)
	if err = services.DeleteJob(db, job); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = h.reloadJobsOnNode(job); err != nil {
		render.Render(w, r, api.ErrNodeResponse(err))
		return
	}
	render.Render(w, r, api.OK("{}"))
}

// RunJob ...
func (h *Handler) RunJob(w http.ResponseWriter, r *http.Request) {
	var err error
	job := r.Context().Value(jobKey).(*models.Job)
	endpoint, err := h.runtime.FindNode(property.NodeIP(job.Node.IP))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.RunJobOnNode(job, endpoint); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.OK("{}"))
}

func (h *Handler) reloadJobsOnNode(job *models.Job) error {
	endpoint, err := h.runtime.FindNode(property.NodeIP(job.Node.IP))
	if err != nil {
		return err
	}
	return services.ReloadJobsOnNode(endpoint)
}
