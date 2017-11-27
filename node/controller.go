package node

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Controller represents a working node
type Controller struct {
	config *config
	db     *gorm.DB
	server *gin.Engine
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
	return &Controller{
		config: cfg,
		db:     db,
		server: server,
	}, nil
}

// Run starts work
func (r *Controller) Run() {
	logrus.Info("node is running...")
	var err error
	if err = r.autoMigrate(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("server is running, port: " + r.config.serverPort())
	r.registerRoutes()
	if err = r.server.Run(":" + r.config.serverPort()); err != nil {
		logrus.Fatal(err)
	}
}

func (r *Controller) registerRoutes() {
	h := handler{
		db: r.db,
	}
	r.server.GET("/k", h.k)
	r.server.POST("/jobs/new", h.jobsNew)
}
