package main

import (
	"github.com/Colm3na/cosmos-opt-api/controller"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/metrics"
	"github.com/Colm3na/cosmos-opt-api/repository"
	"github.com/Colm3na/cosmos-opt-api/service"
)

func main() {
	log := logger.NewLogger()
	metrics := metrics.New()
	repository := repository.NewRepository(log)
	service := service.NewService(repository, log)
	api := controller.NewController(service, metrics, log)
	api.Run()
}
