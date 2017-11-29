package node

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type handler struct {
	db     *gorm.DB
	daemon *daemon
}

func (h *handler) k(c *gin.Context) {
	c.String(http.StatusOK, "Sometimes to love someone, you gotta be a stranger. -- Blade Runner 2049")
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
	if err = h.daemon.reloadCron(); err != nil {
		logrus.Errorf("failed to reload cron jobs, err: %s", err.Error())
	}
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
