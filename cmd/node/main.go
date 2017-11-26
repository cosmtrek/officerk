package main

import (
	"flag"

	"github.com/cosmtrek/officerk/node"
	"github.com/sirupsen/logrus"
)

var cfg string

func init() {
	flag.StringVar(&cfg, "c", "", "config file")
}

func main() {
	flag.Parse()

	var err error
	ctr, err := node.NewController(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	ctr.Run()
}
