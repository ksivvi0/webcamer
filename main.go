package main

import (
	"github.com/sirupsen/logrus"
	"webcamer/config"
	"webcamer/tg_engine"
	"webcamer/webcamer"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		logrus.Error(err)
		return
	}
	engine, err := tg_engine.NewEngine(cfg.Token, webcamer.NewWebcamer(cfg.DefaultDeviceID), cfg.Admins...)
	if err != nil {
		logrus.Error(err)
		return
	}
	engine.Run()

}
