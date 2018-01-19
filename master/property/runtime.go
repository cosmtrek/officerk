package property

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/utils"
)

// NodeIP ...
type NodeIP string

// NodeService ...
type NodeService string

// Runtime saves nodes info, etc.
type Runtime struct {
	IP      string
	Port    string
	IsDebug bool

	sync.RWMutex
	Nodes map[NodeIP]NodeService
}

// NewRuntime ...
func NewRuntime() (*Runtime, error) {
	ip, err := utils.GetIP()
	if err != nil {
		return nil, err
	}
	nodes := make(map[NodeIP]NodeService)
	return &Runtime{
		IP:    ip,
		Nodes: nodes,
	}, nil
}

// AddNode ...
func (r *Runtime) AddNode(ip NodeIP, service NodeService) {
	r.Lock()
	defer r.Unlock()
	r.Nodes[ip] = service
}

// DeleteNode ...
func (r *Runtime) DeleteNode(ip NodeIP) {
	r.Lock()
	defer r.Unlock()
	delete(r.Nodes, ip)
}

// FindNode ...
func (r *Runtime) FindNode(ip NodeIP) (NodeService, error) {
	r.RLock()
	defer r.RUnlock()
	s, ok := r.Nodes[ip]
	if !ok {
		return NodeService("nil"), errors.New("failed to find node")
	}
	return s, nil
}

// Inspect ...
func (r *Runtime) Inspect() {
	logrus.Infof("master ip: %s, port: %s", r.IP, r.Port)
	if r.IsDebug {
		logrus.Debugln("DEBUG mode")
	}
}
