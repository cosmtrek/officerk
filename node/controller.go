package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

// Controller represents a working node
type Controller struct {
	config *config
	db     *gorm.DB
	server *gin.Engine
	daemon *daemon
}

// NewController returns a node
func NewController(path string) (*Controller, error) {
	var err error
	cfg, err := newConfig(path)
	if err != nil {
		return nil, err
	}
	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}
	server := gin.Default()
	daemon, err := newDaemon(db)
	if err != nil {
		return nil, err
	}
	return &Controller{
		config: cfg,
		db:     db,
		server: server,
		daemon: daemon,
	}, nil
}

// Run starts work
func (r *Controller) Run() {
	logrus.Info("node is running...")
	var err error
	if err = r.autoMigrate(); err != nil {
		logrus.Fatal(err)
	}
	go func() {
		logrus.Info("daemon is running...")
		if err = r.daemon.run(); err != nil {
			logrus.Fatal(err)
		}
	}()
	logrus.Info("server is running, port: " + r.config.serverPort())
	r.registerRoutes()
	if err = r.server.Run(":" + r.config.serverPort()); err != nil {
		logrus.Fatal(err)
	}
}

func (r *Controller) registerRoutes() {
	h := handler{
		db:     r.db,
		daemon: r.daemon,
	}
	r.server.GET("/k", h.k)

	v1 := r.server.Group("/v1")
	jobs := v1.Group("/jobs")
	jobs.GET("/", h.jobsIndex)
	jobs.POST("/new", h.jobsNew)
	jobs.GET("/:id", h.jobsDetail)
	jobs.PUT("/:id", h.jobsUpdate)
	jobs.DELETE("/:id", h.jobsDelete)

	admin := r.server.Group("/admin")
	admin.GET("/jobs/run/:id", h.jobsRun)
	admin.GET("/jobs/reload", h.jobsReload)
}
