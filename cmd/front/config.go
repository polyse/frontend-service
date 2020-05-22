package main

import (
	"time"

	"github.com/polyse/frontend-service/internal/api"
	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env"
)

// Config is main application configuration structure.
type config struct {
	Listen   string        `env:"LISTEN" envDefault:"localhost:9900"`
	DB       string        `env:"DB" envDefault:"localhost:9000"`
	Timeout  time.Duration `env:"TIMEOUT" envDefault:"10ms"`
	LogLevel string        `env:"LOG_LEVEL" envDefault:"info"`
	LogFmt   string        `env:"LOG_FMT" envDefault:"console"`
}

func load() (*config, error) {
	log.Debug().Msg("loading configuration")
	cfg := &config{}

	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func initFrontendServiceCfg(c *config) (api.AppConfig, error) {
	return api.AppConfig{Timeout: c.Timeout, NetInterface: c.Listen, DB: c.DB}, nil
}
