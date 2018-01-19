package master

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/master/api"
	"github.com/cosmtrek/officerk/master/property"
	"github.com/cosmtrek/officerk/utils"
)

var (
	etcdNodeDir = "/officerk-nodes/"
)

// Controller represents a working node
type Controller struct {
	runtime *property.Runtime
	config  *utils.Config
	db      *gorm.DB
	router  *chi.Mux
	etcd    *clientv3.Client
}

// NewController returns a node
func NewController(cfg *utils.Config, opt *Option) (*Controller, error) {
	var err error

	runtime, err := property.NewRuntime()
	if err != nil {
		return nil, err
	}
	if len(opt.Port) > 0 {
		runtime.Port = opt.Port
	} else {
		runtime.Port = cfg.MasterServerPort()
	}
	runtime.IsDebug = opt.IsDebug

	etcd, err := utils.NewEtcdClient(cfg)
	if err != nil {
		return nil, err
	}

	db, err := utils.OpenDB(cfg)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	return &Controller{
		runtime: runtime,
		config:  cfg,
		db:      db,
		router:  router,
		etcd:    etcd,
	}, nil
}

// Run starts work
func (ctr *Controller) Run() {
	ctr.runtime.Inspect()

	var err error
	if err = ctr.autoMigrate(); err != nil {
		logrus.Fatal(err)
	}

	go func() {
		// TODO: signal
		ctr.watchNodes()
	}()

	ctr.registerRoutes()
	if err = http.ListenAndServe(":"+ctr.runtime.Port, ctr.router); err != nil {
		logrus.Fatal(err)
	}
}

func (ctr *Controller) registerRoutes() {
	h := api.NewHandler(ctr.db, ctr.runtime)
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
			r.Get("/logs", h.GetJobLogs)
		})
	})

	r.Get("/joblogs/{joblogID}/", h.GetJobLog)

	r.Route("/nodes", func(r chi.Router) {
		r.Get("/", h.ListNodes)
		r.Post("/", h.CreateNode)
	})
}

func (ctr *Controller) watchNodes() {
	runtime := ctr.runtime
	kv := clientv3.NewKV(ctr.etcd)
	curRevision := int64(0)

	logrus.Debug("searching nodes...")
	for {
		rangeResp, err := kv.Get(context.TODO(), etcdNodeDir, clientv3.WithPrefix())
		if err != nil {
			continue
		}

		for _, kv := range rangeResp.Kvs {
			v := strings.Split(string(kv.Value), ":")
			runtime.AddNode(property.NodeIP(v[0]), property.NodeService(kv.Value))
		}

		// find the latest revision to watch
		curRevision = rangeResp.Header.Revision + 1
		break
	}

	watcher := clientv3.NewWatcher(ctr.etcd)
	watchChan := watcher.Watch(context.TODO(), etcdNodeDir, clientv3.WithPrefix(), clientv3.WithRev(curRevision))
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				key := string(event.Kv.Key)
				value := string(event.Kv.Value)
				logrus.Debugf("PUT event: %s - %s", key, value)
				v := strings.Split(key, "/")
				runtime.AddNode(property.NodeIP(v[2]), property.NodeService(value))
			case mvccpb.DELETE:
				key := string(event.Kv.Key)
				logrus.Debugf("DELETE event: %s", key)
				v := strings.Split(key, "/")
				runtime.DeleteNode(property.NodeIP(v[2]))
			}
			logrus.Debugf("nodes: %+v", runtime.Nodes)
		}
	}
}
