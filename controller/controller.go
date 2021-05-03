package controller

import (
	"net/http"

	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/models"
	"github.com/Colm3na/cosmos-opt-api/service"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	service service.Service
	logger  *logger.Logger
}

func NewController(service service.Service, logger *logger.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c *Controller) GetBlocktime(e echo.Context) error {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), e)

	latestBlockHeigh, latestBlockTime, err := c.service.GetBlockHeighAndTime(-1)
	if err != nil {
		c.logger.LogOnErrorWithContext(c.logger.GetContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	avg, err := c.service.GetAverageBlockTime(latestBlockHeigh, latestBlockTime)
	if err != nil {
		c.logger.LogOnErrorWithContext(c.logger.GetContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	blocktime := models.BlockTime{Average: avg}

	c.logger.LogOnExitWithContext(c.logger.GetContext(), blocktime)
	return e.JSON(http.StatusOK, blocktime)
}

func (c *Controller) GetValidatorUptime(e echo.Context) error {
	c.logger.LogOnEntryWithContext(c.logger.GetContext(), e)

	cosmosvaloper := e.Param("id")
	uptime, err := c.service.GetValidatorUptime(cosmosvaloper)
	if err != nil {
		c.logger.LogOnErrorWithContext(c.logger.GetContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	c.logger.LogOnExitWithContext(c.logger.GetContext(), uptime)
	return e.JSON(http.StatusOK, uptime)
}
