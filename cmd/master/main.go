package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/cosmtrek/officerk/master"
	"github.com/cosmtrek/officerk/utils"
)

var (
	BuildTimestamp string
	Version        string

	cfgPath string
	isDebug bool
	port    string
)

func init() {
	flag.StringVar(&cfgPath, "c", "", "config file")
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&port, "port", "", "master port")
	flag.Parse()

	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	fmt.Println(utils.MasterLogo(BuildTimestamp, Version))

	var err error
	cfg, err := utils.NewConfig(cfgPath)
	if err != nil {
		logrus.Fatal(err)
	}

	opt := &master.Option{
		Port:    port,
		IsDebug: isDebug,
	}
	ctr, err := master.NewController(cfg, opt)
	if err != nil {
		logrus.Fatal(err)
	}
	ctr.Run()
}
