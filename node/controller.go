package node

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Controller represents a working node
type Controller struct {
	config *config
	db     *gorm.DB
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
	return &Controller{
		config: cfg,
		db:     db,
	}, nil
}

// Run starts work
func (r *Controller) Run() {
	logrus.Info("node is running...")
	var err error
	if err = r.autoMigrate(); err != nil {
		logrus.Fatal(err)
	}
}
