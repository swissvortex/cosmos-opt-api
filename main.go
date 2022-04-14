package main

import (
	"github.com/swissvortex/cosmos-opt-api/controller"
	"github.com/swissvortex/cosmos-opt-api/logger"
	"github.com/swissvortex/cosmos-opt-api/metrics"
	"github.com/swissvortex/cosmos-opt-api/repository"
	"github.com/swissvortex/cosmos-opt-api/service"
)

func main() {
	log := logger.NewLogger()
	metrics := metrics.New()
	repository := repository.NewRepository(log)
	service := service.NewService(repository, log)
	api := controller.NewController(service, metrics, log)
	api.Run()
}
