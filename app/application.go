package app

import (
	"github.com/ApplyLogic/mqtt-subscriber/config"
	"github.com/ApplyLogic/mqtt-subscriber/internal/middleware"
	mqttLib "github.com/ApplyLogic/mqtt-subscriber/mqtt"
)

type App struct {
	Config     *config.Config
	Subscriber *mqttLib.Subscriber
	Logger     *middleware.Logger
}

func (a *App) Initialize(cfg *config.Config) {
	// Initialize Logger
	a.Config = cfg
	a.Logger = &middleware.Logger{}
	a.Logger.Initialize(cfg)
	a.Subscriber = mqttLib.New(a.Config, a.Logger)
}

func (a *App) Start() {
	a.Subscriber.Start()
}
