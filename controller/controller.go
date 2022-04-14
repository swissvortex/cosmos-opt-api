package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swissvortex/cosmos-opt-api/constants"
	"github.com/swissvortex/cosmos-opt-api/logger"
	"github.com/swissvortex/cosmos-opt-api/metrics"
	"github.com/swissvortex/cosmos-opt-api/models"
	"github.com/swissvortex/cosmos-opt-api/service"
)

type Controller interface {
	Run()
}

type controller struct {
	service service.Service
	metric  metrics.Metric
	log     logger.Logger
	*echo.Echo
}

func NewController(service service.Service, metric metrics.Metric, log logger.Logger) Controller {
	return &controller{
		service: service,
		metric:  metric,
		log:     log,
		Echo:    echo.New(),
	}
}

func (c *controller) Run() {
	c.log.EntryWithContext(c.log.FileContext())
	c.ConfigureMiddleware()
	c.ConfigureEndpoints()
	c.InitMetrics()
	go c.UpdateBlocktime()
	c.log.ErrorWithContext(c.log.FileContext(), c.Echo.Start(constants.ServerHost+":"+constants.ServerPort))
}

func (c *controller) ConfigureMiddleware() {
	c.log.EntryWithContext(c.log.FileContext())

	c.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{constants.ServerHost},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	c.Use(middleware.Logger())
	c.Use(middleware.Recover())
	prommetheus := prometheus.NewPrometheus(constants.ProjectName, nil)
	prommetheus.Use(c.Echo)

	c.log.ExitWithContext(c.log.FileContext())
}

func (c *controller) ConfigureEndpoints() {
	health := c.Group("")
	{
		health.GET("/blocktime", c.GetBlocktimeApi)
		health.GET("/validator/:cosmosvaloper", c.GetValidatorUptime)
	}
}

func (c *controller) InitMetrics() {
	c.metric.NewGauge(constants.AverageBlocktimeGaugeName, constants.AverageBlocktimeGaugeHelp)
	c.metric.NewGauge(constants.ValidatorUptimeGaugeName, constants.ValidatorUptimeGaugeHelp)
}

func (c *controller) UpdateBlocktime() {
	for {
		time.Sleep(time.Duration(constants.PrometheusUpdateTime) * time.Second)
		blocktime, err := c.GetBlocktime()
		if err != nil {
			c.log.ErrorWithContext(c.log.FileContext(), err)
		}
		c.log.DebugWithContext(c.log.FileContext(), blocktime)
	}
}

func (c *controller) GetBlocktime() (*models.BlockTime, error) {
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

	c.metric.SetGauge(constants.AverageBlocktimeGaugeName, avg)
	blocktime := models.BlockTime{Average: avg}

	c.log.ExitWithContext(c.log.FileContext(), blocktime)
	return &blocktime, nil
}

func (c *controller) GetBlocktimeApi(e echo.Context) error {
	c.log.EntryWithContext(c.log.FileContext(), e)

	blocktime, err := c.GetBlocktime()
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	c.log.ExitWithContext(c.log.FileContext(), blocktime)
	return e.JSON(http.StatusOK, blocktime)
}

func (c *controller) GetValidatorUptime(e echo.Context) error {
	c.log.EntryWithContext(c.log.FileContext(), e)

	cosmosvaloper := e.Param("cosmosvaloper")
	uptime, err := c.service.GetValidatorUptime(cosmosvaloper)
	if err != nil {
		c.log.ErrorWithContext(c.log.FileContext(), err)
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	c.metric.SetGauge(constants.ValidatorUptimeGaugeName, float64((100-uptime.MissedBlocks)*100/uptime.OverBlocks))

	c.log.ExitWithContext(c.log.FileContext(), uptime)
	return e.JSON(http.StatusOK, uptime)
}
