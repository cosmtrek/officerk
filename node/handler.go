package node

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type handler struct {
	db     *gorm.DB
	daemon *daemon
}

func (h *handler) k(c *gin.Context) {
	c.String(http.StatusOK, "Sometimes to love someone, you gotta be a stranger. -- Blade Runner 2049")
}

func (h *handler) jobsIndex(c *gin.Context) {
	var err error
	data, err := h.daemon.getJobs()
	if err != nil {
		responseBadRequest(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
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
	h.daemon.restartCron()
	responseOK(c)
}

func (h *handler) jobsRun(c *gin.Context) {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		responseBadRequest(c, err)
		return
	}
	err = h.daemon.runJob(uint(id))
	if err != nil {
		responseBadRequest(c, err)
		return
	}
	responseOK(c)
}

func (h *handler) jobsDetail(c *gin.Context) {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		responseBadRequest(c, err)
		return
	}
	data, err := h.daemon.getJob(uint(id))
	if err != nil {
		responseBadRequest(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *handler) jobsUpdate(c *gin.Context) {
	var err error
	var jr JobRequest
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		responseBadRequest(c, err)
		return
	}
	if err = c.BindJSON(&jr); err != nil {
		responseBadRequest(c, err)
		return
	}
	if err = jr.validate(); err != nil {
		responseBadRequest(c, err)
		return
	}
	var tasks []Task
	for _, task := range jr.Tasks {
		t := Task{
			JobID:     uint(id),
			Name:      task.Name,
			Command:   task.Command,
			NextTasks: task.NextTasks,
		}
		if len(task.ID) > 0 {
			id, err := strconv.Atoi(task.ID)
			if err != nil || id <= 0 {
				responseBadRequest(c, err)
				return
			}
			t.ID = uint(id)
		}
		tasks = append(tasks, t)
	}
	job := Job{
		Name:      jr.Name,
		Schedule:  jr.Schedule,
		RoutePath: jr.RoutePath,
		Tasks:     tasks,
	}
	if err = h.daemon.updateJob(uint(id), job); err != nil {
		responseBadRequest(c, err)
		return
	}
	h.daemon.restartCron()
	responseOK(c)
}

func (h *handler) jobsDelete(c *gin.Context) {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		responseBadRequest(c, err)
		return
	}
	err = h.daemon.deleteJob(uint(id))
	if err != nil {
		responseBadRequest(c, err)
		return
	}
	h.daemon.restartCron()
	responseOK(c)
}

func (h *handler) jobsReload(c *gin.Context) {
	h.daemon.restartCron()
	responseOK(c)
}

func responseOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func responseBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
