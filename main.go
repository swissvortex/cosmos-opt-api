package main

import (
	"log"
	"time"

	"github.com/Colm3na/cosmos-opt-api/controller"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/repository"
	"github.com/Colm3na/cosmos-opt-api/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	. "github.com/Colm3na/cosmos-opt-api/constants"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	logger := logger.NewLogger()
	logger.SetLoggingLevel(LOGGING_LEVEL)
	repository := repository.NewRepository(logger)
	service := service.NewService(repository, logger)
	api := controller.NewController(service, logger)

	b, err := tb.NewBot(tb.Settings{
		Token:  "1707302307:AAHAGVGMFs28GW1_2m8ND938X_Yu3jMYUOE",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hola jefe")
	})

	b.Start()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/blocktime", api.GetBlocktime)
	e.GET("/validator/:id", api.GetValidatorUptime)
	e.Logger.Fatal(e.Start(":47567"))
}
