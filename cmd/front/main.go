package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/polyse/frontend-service/internal/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xlab/closer"
)

func main() {
	defer closer.Close()

	cfg, err := load()
	if err != nil {
		log.Err(err).Msg("error while loading config")
		return
	}

	if err = initLogger(cfg); err != nil {
		log.Err(err).Msg("error while configure logger")
		return
	}
	log.Debug().Msg("logger initialized")

	log.Debug().Msg("starting di container")
	a, cancel, err := initFrontendService(cfg)
	if err != nil {
		log.Err(err).Msg("error while init wire")
		return
	}
	closer.Bind(cancel)

	log.Debug().Msg("starting frontend web application")
	if err = a.Run(); err != nil {
		log.Err(err).Msg("error while starting api app")
	}
}

func initFrontendServiceCfg(c *config) (api.AppConfig, error) {
	return api.AppConfig{Timeout: c.Timeout, NetInterface: c.Listen}, nil
}

func initLogger(c *config) error {
	log.Debug().Msg("initialize logger")
	logLvl, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLvl)
	switch c.LogFmt {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	case "json":
	default:
		return fmt.Errorf("unknown output format %s", c.LogFmt)
	}
	return nil
}
