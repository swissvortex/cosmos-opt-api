package controller

import (
	"net/http"

	"github.com/Colm3na/cosmos-opt-api/constants"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/metrics"
	"github.com/Colm3na/cosmos-opt-api/models"
	"github.com/Colm3na/cosmos-opt-api/service"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	service service.Service
	metric  metrics.Metric
	log     logger.Logger
}

func NewController(service service.Service, metric metrics.Metric, log logger.Logger) *Controller {
	InitMetrics(metric)
	return &Controller{
		service: service,
		metric:  metric,
		log:     log,
	}
}

func InitMetrics(metric metrics.Metric) {
	metric.NewGauge("blocktime_average", "Average block time in seconds")
	metric.NewGauge("validator_uptime", "Validator uptime in seconds")
}

func (c *Controller) GetBlocktime() (*models.BlockTime, error) {
	c.log.EntryWithContext(c.log.FileContext())

	latestBlockHeigh, latestBlockTime, err := c.service.GetBlockHeighAndTime(constants.LatestBlock)
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return nil, err
	}

	avg, err := c.service.GetAverageBlockTime(latestBlockHeigh, latestBlockTime)
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return nil, err
	}

	c.metric.SetGauge("blocktime_average", avg)
	blocktime := models.BlockTime{Average: avg}

	c.log.ExitWithContext(c.log.FileContext(), blocktime)
	return &blocktime, nil
}

func (c *Controller) GetBlocktimeApi(e echo.Context) error {
	c.log.EntryWithContext(c.log.FileContext(), e)

	blocktime, err := c.GetBlocktime()
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	c.log.ExitWithContext(c.log.FileContext(), blocktime)
	return e.JSON(http.StatusOK, blocktime)
}

func (c *Controller) GetValidatorUptime(e echo.Context) error {
	c.log.EntryWithContext(c.log.FileContext(), e)

	cosmosvaloper := e.Param("id")
	uptime, err := c.service.GetValidatorUptime(cosmosvaloper)
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	c.metric.SetGauge("validator_uptime", float64((100-uptime.MissedBlocks)*100/uptime.OverBlocks))

	c.log.ExitWithContext(c.log.FileContext(), uptime)
	return e.JSON(http.StatusOK, uptime)
}
