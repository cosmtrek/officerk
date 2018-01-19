package api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/cosmtrek/officerk/models"
)

// TaskLogRequest ...
type TaskLogRequest struct {
	*models.TaskLog
}

// Bind for post-processing TaskLogRequest
func (t *TaskLogRequest) Bind(r *http.Request) error {
	return nil
}

// TaskLogResponse ...
type TaskLogResponse struct {
	*models.TaskLog
}

// NewTaskLogResponse ...
func NewTaskLogResponse(jobLog *models.TaskLog) *TaskLogResponse {
	return &TaskLogResponse{TaskLog: jobLog}
}

// Render for TaskLogResponse
func (jr TaskLogResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TaskLogListResponse ...
type TaskLogListResponse []*TaskLogResponse

// NewTaskLogListResponse ...
func NewTaskLogListResponse(logs []*models.TaskLog) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, log := range logs {
		list = append(list, NewTaskLogResponse(log))
	}
	return list
}
