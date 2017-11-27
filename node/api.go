package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	Name      string `json:"name"`
	Command   string `json:"command"`
	NextTasks string `json:"next_tasks"` // "task1,task2,task3"
}

// TODO
func (j *JobRequest) validate() error {
	return nil
}

// TODO
func (j *JobRequest) save(db *gorm.DB) error {
	return nil
}

type handler struct {
	db *gorm.DB
}

func (h *handler) jobsNew(c *gin.Context) {
	var err error
	var jr JobRequest
	if err = c.BindJSON(&jr); err != nil {
		responseBadRequest(c, err)
		return
	}
	if err = jr.validate(); err != nil {
		responseBadRequest(c, err)
		return
	}
	if err = jr.save(h.db); err != nil {
		responseBadRequest(c, err)
		return
	}
	responseOK(c)
}

func (h *handler) k(c *gin.Context) {
	c.String(http.StatusOK, "Sometimes to love someone, you gotta be a stranger. -- Blade Runner 2049")
}

func responseOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func responseBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}
