package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/cosmtrek/officerk/node/daemon"
	"github.com/cosmtrek/officerk/utils/api"
)

// Handler ...
type Handler struct {
	jobDaemon *daemon.Engine
}

// NewHandler ...
func NewHandler(dm *daemon.Engine) *Handler {
	return &Handler{
		jobDaemon: dm,
	}
}

// ReloadJobs ...
func (h *Handler) ReloadJobs(w http.ResponseWriter, r *http.Request) {
	h.jobDaemon.Reload()
	render.Render(w, r, api.OK("{}"))
}

// RunJob ...
func (h *Handler) RunJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobID")
	if jobID == "" {
		render.Render(w, r, api.ErrNotFound)
		return
	}
	var err error
	id, err := strconv.Atoi(jobID)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = h.jobDaemon.RunJob(uint(id)); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, api.OK("{}"))
}
