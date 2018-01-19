package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/node"
	"github.com/cosmtrek/officerk/utils"
)

var (
	cfgPath string
	isDebug bool
	ip      string
	port    string
)

func init() {
	flag.StringVar(&cfgPath, "c", "", "config file")
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&ip, "ip", "", "node ip")
	flag.StringVar(&port, "port", "", "node port")
	flag.Parse()

	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	fmt.Println(utils.NodeLogo())

	var err error
	cfg, err := utils.NewConfig(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}
	opt := &node.Option{
		IP:      ip,
		Port:    port,
		IsDebug: isDebug,
	}
	ctr, err := node.NewController(cfg, opt)
	if err != nil {
		logrus.Fatal(err)
	}
	ctr.Run()
}
