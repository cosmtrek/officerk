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
	nodes map[NodeIP]NodeService
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
		nodes: nodes,
	}, nil
}

// AddNode ...
func (r *Runtime) AddNode(ip NodeIP, service NodeService) {
	r.Lock()
	defer r.Unlock()
	r.nodes[ip] = service
}

// DeleteNode ...
func (r *Runtime) DeleteNode(ip NodeIP) {
	r.Lock()
	defer r.Unlock()
	delete(r.nodes, ip)
}

// FindNode ...
func (r *Runtime) FindNode(ip NodeIP) (NodeService, error) {
	r.RLock()
	defer r.RUnlock()
	s, ok := r.nodes[ip]
	if !ok {
		logrus.Debugf("not found node ip: %s", string(ip))
		return NodeService("nil"), errors.New("failed to find the node, is it online?")
	}
	return s, nil
}

// Nodes ...
func (r *Runtime) Nodes() []string {
	r.RLock()
	defer r.RUnlock()
	nodes := make([]string, 0)
	for k, _ := range r.nodes {
		nodes = append(nodes, string(k))
	}
	return nodes
}

// IsOnline ...
func (r *Runtime) IsOnline(ip string) bool {
	r.RLock()
	defer r.RUnlock()
	for k, _ := range r.nodes {
		if ip == string(k) {
			return true
		}
	}
	return false
}

// Inspect ...
func (r *Runtime) Inspect() {
	logrus.Infof("master ip: %s, port: %s", r.IP, r.Port)
	if r.IsDebug {
		logrus.Debugln("DEBUG mode")
	}
}
