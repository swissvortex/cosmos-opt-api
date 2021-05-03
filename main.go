package main

import (
	"github.com/Colm3na/cosmos-opt-api/controller"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/repository"
	"github.com/Colm3na/cosmos-opt-api/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	. "github.com/Colm3na/cosmos-opt-api/constants"
)

func main() {
	logger := logger.NewLogger()
	logger.SetLoggingLevel(LOGGING_LEVEL)
	cryptoRepository := repository.NewRepository(logger)
	service := service.NewService(cryptoRepository, logger)
	api := controller.NewController(service, logger)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/blocktime", api.GetBlocktime)
	e.GET("/validator/:id", api.GetValidatorUptime)
	e.Logger.Fatal(e.Start(":47567"))
}
