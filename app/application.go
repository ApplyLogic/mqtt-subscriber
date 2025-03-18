package app

import (
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	"github.com/ApplyLogic/mqtt-subscriber/internal/middleware"
	mqttLib "github.com/ApplyLogic/mqtt-subscriber/mqtt"
	"os"
	"os/signal"
	"syscall"
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
	if token := a.Subscriber.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	if token := a.Subscriber.Client.Subscribe(a.Config.MQTT.Topic, 2, mqttLib.MessageSubHandler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	}

	fmt.Println("Subscribed to topic:", a.Config.MQTT.Topic)

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	a.Subscriber.Client.Unsubscribe(a.Config.MQTT.Topic)
	a.Subscriber.Client.Disconnect(250)
}
