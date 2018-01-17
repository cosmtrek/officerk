package master

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/master/api"
	"github.com/cosmtrek/officerk/utils"
)

// Controller represents a working node
type Controller struct {
	config *utils.Config
	db     *gorm.DB
	router *chi.Mux
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
	return &Controller{
		config: cfg,
		db:     db,
		router: router,
	}, nil
}

// Run starts work
func (ctr *Controller) Run() {
	var err error
	if err = ctr.autoMigrate(); err != nil {
		logrus.Fatal(err)
	}

	ctr.registerRoutes()
	logrus.Info("server is running, port: " + ctr.config.MasterServerPort())
	if err = http.ListenAndServe(":"+ctr.config.MasterServerPort(), ctr.router); err != nil {
		logrus.Fatal(err)
	}
}

func (ctr *Controller) registerRoutes() {
	h := api.NewHandler(ctr.db)
	r := ctr.router
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/k", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sometimes to love someone, you gotta be a stranger. -- Blade Runner 2049"))
	})

	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", h.ListJobs)
		r.Post("/", h.CreateJob)

		r.Route("/{jobID}", func(r chi.Router) {
			r.Use(api.JobCtx)
			r.Get("/", h.GetJob)
			r.Put("/", h.UpdateJob)
			r.Delete("/", h.DeleteJob)
			r.Get("/run", h.RunJob)
		})
	})

	// TODO: nodes
	r.Route("/nodes", func(r chi.Router) {
		r.Get("/", h.ListNodes)
		r.Post("/", h.CreateNode)
	})
}
