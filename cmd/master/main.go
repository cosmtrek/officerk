package main

import (
	"flag"

	"github.com/cosmtrek/officerk/master"
	"github.com/sirupsen/logrus"
)

var cfg string

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	flag.StringVar(&cfg, "c", "", "config file")
}

func main() {
	flag.Parse()

	var err error
	ctr, err := master.NewController(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	ctr.Run()
}
