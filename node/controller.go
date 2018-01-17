package node

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/node/api"
	"github.com/cosmtrek/officerk/node/daemon"
	"github.com/cosmtrek/officerk/utils"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

// Controller represents a working node
type Controller struct {
	config    *utils.Config
	db        *gorm.DB
	router    *chi.Mux
	jobDaemon *daemon.Engine
}

// NewController returns a node
func NewController(path string) (*Controller, error) {
	var err error
	cfg, err := utils.NewConfig(path)
	if err != nil {
		return nil, err
	}
	db, err := utils.OpenDB(cfg)
	if err != nil {
		return nil, err
	}
	router := chi.NewRouter()
	jobDaemon, err := daemon.NewEngine(db)
	if err != nil {
		return nil, err
	}
	return &Controller{
		config:    cfg,
		db:        db,
		router:    router,
		jobDaemon: jobDaemon,
	}, nil
}

// Run starts work
func (ctr *Controller) Run() {
	logrus.Info("node is running...")
	var err error
	go func() {
		// TODO: signal
		logrus.Info("job daemon is running...")
		if err = ctr.jobDaemon.Run(); err != nil {
			logrus.Fatal(err)
		}
	}()
	ctr.registerRoutes()
	logrus.Info("server is running, port: " + ctr.config.NodeServerPort())
	if err = http.ListenAndServe(":"+ctr.config.NodeServerPort(), ctr.router); err != nil {
		logrus.Fatal(err)
	}
}

func (ctr *Controller) registerRoutes() {
	h := api.NewHandler(ctr.jobDaemon)
	r := ctr.router
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/reload_jobs", h.ReloadJobs)
	r.Get("/jobs/{jobID}/run", h.RunJob)
}
