package utils

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

// NewEtcdClient connects etcd
func NewEtcdClient(c *Config) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints(),
		DialTimeout: 5 * time.Second,
	})
}
