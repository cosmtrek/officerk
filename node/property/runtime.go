package property

import (
	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/utils"
)

// Runtime ...
type Runtime struct {
	IP      string
	Port    string
	IsDebug bool
}

// NewRuntime ...
func NewRuntime() (*Runtime, error) {
	ip, err := utils.GetIP()
	if err != nil {
		return nil, err
	}
	return &Runtime{
		IP: ip,
	}, nil
}

// Inspect ...
func (r *Runtime) Inspect() {
	logrus.Infof("node ip: %s, port: %s", r.IP, r.Port)
	if r.IsDebug {
		logrus.Debugln("DEBUG mode")
	}
}
