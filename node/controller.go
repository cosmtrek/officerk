package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/node/api"
	"github.com/cosmtrek/officerk/node/daemon"
	"github.com/cosmtrek/officerk/node/property"
	"github.com/cosmtrek/officerk/utils"
)

var (
	etcdNodeDir = "/officerk-nodes/"
)

// Controller represents a working node
type Controller struct {
	runtime   *property.Runtime
	config    *utils.Config
	db        *gorm.DB
	router    *chi.Mux
	jobDaemon *daemon.Engine
	etcd      *clientv3.Client
}

// NewController returns a node
func NewController(cfg *utils.Config, opt *Option) (*Controller, error) {
	var err error

	runtime, err := property.NewRuntime()
	if err != nil {
		return nil, err
	}
	if len(opt.IP) > 0 {
		runtime.IP = opt.IP
	}
	if len(opt.Port) > 0 {
		runtime.Port = opt.Port
	} else {
		runtime.Port = cfg.NodeServerPort()
	}
	runtime.IsDebug = opt.IsDebug

	db, err := utils.OpenDB(cfg)
	if err != nil {
		return nil, err
	}
	etcd, err := utils.NewEtcdClient(cfg)
	if err != nil {
		return nil, err
	}
	router := chi.NewRouter()
	jobDaemon, err := daemon.NewEngine(db, runtime)
	if err != nil {
		return nil, err
	}
	return &Controller{
		config:    cfg,
		db:        db,
		router:    router,
		jobDaemon: jobDaemon,
		etcd:      etcd,
		runtime:   runtime,
	}, nil
}

// Run starts work
func (ctr *Controller) Run() {
	ctr.runtime.Inspect()

	var err error

	go func() {
		// TODO: signal
		ctr.keepAlive()
	}()
	go func() {
		// TODO: signal
		logrus.Info("job daemon is running...")
		if err = ctr.jobDaemon.Run(); err != nil {
			logrus.Fatal(err)
		}
	}()

	ctr.registerRoutes()
	if err = http.ListenAndServe(":"+ctr.runtime.Port, ctr.router); err != nil {
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

func (ctr *Controller) keepAlive() {
	var curLeaseID clientv3.LeaseID
	kv := clientv3.NewKV(ctr.etcd)
	lease := clientv3.NewLease(ctr.etcd)
	value := fmt.Sprintf("%s:%s", ctr.runtime.IP, ctr.runtime.Port)

	logrus.Debug("calling master...")
	for {
		if curLeaseID == 0 {
			leaseResp, err := lease.Grant(context.TODO(), 10)
			if err != nil {
				goto SLEEP
			}

			key := fmt.Sprintf("%s%s", etcdNodeDir, ctr.runtime.IP)
			if _, err := kv.Put(context.TODO(), key, value, clientv3.WithLease(leaseResp.ID)); err != nil {
				goto SLEEP
			}
			curLeaseID = leaseResp.ID
		} else {
			logrus.Debugf("keepalive curLeaseID=%d", curLeaseID)
			if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseID); err == rpctypes.ErrLeaseNotFound {
				curLeaseID = 0
				continue
			}
		}
	SLEEP:
		time.Sleep(time.Duration(5) * time.Second)
	}
}
