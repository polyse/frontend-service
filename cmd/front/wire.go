//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/polyse/frontend-service/internal/api"
)

func initFrontendService(c *config) (*api.API, func(), error) {
	wire.Build(api.NewApp, initFrontendServiceCfg)
	return nil, nil, nil
}
