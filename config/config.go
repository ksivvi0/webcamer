package config

import (
	"github.com/codingconcepts/env"
)

type config struct {
	Token           string  `env:"TOKEN" required:"true"`
	Admins          []int64 `env:"ADMINS" delimiter:","`
	DefaultDeviceID int     `env:"DEFAULT_DEVICE_ID" default:"0"`
}

func Init() (*config, error) {
	Cfg := config{}
	if err := env.Set(&Cfg); err != nil {
		return nil, err
	}
	return &Cfg, nil
}
