package node

import (
	"fmt"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
)

const (
	defaultServerPort = "9392"
)

type config struct {
	Master   master   `toml:"master"`
	Database database `toml:"database"`
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

func newConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := config{}
	if err = toml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// MySQL dsn
func (c *config) MySQL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Dbname)
}

func (c *config) serverPort() string {
	port := c.Master.ServerPort
	if port == "" {
		return defaultServerPort
	}
	return port
}

func connectDB(c *config) (*gorm.DB, error) {
	return gorm.Open("mysql", c.MySQL())
}
