package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

const (
	defaultMasterServerPort = "9392"
	defaultNodeServerPort   = "9100"
)

type Config struct {
	Database database `toml:"database"`
	Master   master   `toml:"master"`
	Node     node     `toml:"node"`
}

type database struct {
	Host     string `toml:"localhost"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Dbname   string `toml:"dbname"`
	Password string `toml:"password"`
}

type master struct {
	ServerPort string `toml:"server_port"`
}

type node struct {
	ServerPort string `toml:"server_port"`
}

func NewConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err = toml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// MySQL dsn
func (c *Config) MySQL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Dbname)
}

func (c *Config) MasterServerPort() string {
	port := c.Master.ServerPort
	if port == "" {
		return defaultMasterServerPort
	}
	return port
}

func (c *Config) NodeServerPort() string {
	port := c.Node.ServerPort
	if port == "" {
		return defaultNodeServerPort
	}
	return port
}
