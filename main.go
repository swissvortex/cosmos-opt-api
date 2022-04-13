package main

import (
	"fmt"
	"time"

	"github.com/Colm3na/cosmos-opt-api/controller"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/metrics"
	"github.com/Colm3na/cosmos-opt-api/repository"
	"github.com/Colm3na/cosmos-opt-api/service"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log := logger.NewLogger()
	metrics := metrics.New()
	repository := repository.NewRepository(log)
	service := service.NewService(repository, log)
	api := controller.NewController(service, metrics, log)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/blocktime", api.GetBlocktimeApi)
	e.GET("/validator/:id", api.GetValidatorUptime)
	prommetheus := prometheus.NewPrometheus("cosmos-opt-api", nil)
	prommetheus.Use(e)

	go UpdateBlocktime(api)
	e.Logger.Fatal(e.Start(":8080"))

}

func UpdateBlocktime(api *controller.Controller) {
	for {
		time.Sleep(5 * time.Second)
		blocktime, _ := api.GetBlocktime()
		fmt.Println("Blocktime:", blocktime.Average)
	}
}
