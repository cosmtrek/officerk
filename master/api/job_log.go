package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/services"
	"github.com/cosmtrek/officerk/utils/api"
)

// JobLogRequest ...
type JobLogRequest struct {
	*models.JobLog
}

// Bind for post-processing JobLogRequest
func (j *JobLogRequest) Bind(r *http.Request) error {
	return nil
}

// JobLogResponse ...
type JobLogResponse struct {
	*models.JobLog
}

// NewJobLogResponse ...
func NewJobLogResponse(jobLog *models.JobLog) *JobLogResponse {
	return &JobLogResponse{JobLog: jobLog}
}

// Render for JobLogResponse
func (jr JobLogResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// JobLogListResponse ...
type JobLogListResponse []*JobLogResponse

// NewJobLogListResponse ...
func NewJobLogListResponse(logs []*models.JobLog) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, log := range logs {
		list = append(list, NewJobLogResponse(log))
	}
	return list
}

// GetJobLog ...
func (h *Handler) GetJobLog(w http.ResponseWriter, r *http.Request) {
	var err error
	id := chi.URLParam(r, "joblogID")
	if id == "" {
		render.Render(w, r, api.ErrNotFound)
		return
	}
	joblog := new(models.JobLog)
	err = services.GetJobLog(db, id, joblog)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.OK(NewJobLogResponse(joblog)))
}

// GetJobLogs ...
func (h *Handler) GetJobLogs(w http.ResponseWriter, r *http.Request) {
	var err error
	job := r.Context().Value(jobKey).(*models.Job)
	logs, err := services.GetJobLogs(db, job)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.OK(NewJobLogListResponse(logs)))
}

// ListJobLogs ...
func (h *Handler) ListJobLogs(w http.ResponseWriter, r *http.Request) {
	var err error
	logs, err := services.ListJobLogs(db)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.OK(NewJobLogListResponse(logs)))
}
