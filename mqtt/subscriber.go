package mqtt

import (
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	"github.com/ApplyLogic/mqtt-subscriber/internal/middleware"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
)

type Subscriber struct {
	Client mqtt.Client
	Logger *middleware.Logger
	Config *config.Config
}

func New(cfg *config.Config, logger *middleware.Logger) *Subscriber {

	//tlsConfig, err := configureTLS(cfg)
	//if err != nil {
	//	log.Fatalf("Error configuring TLS: %v", err)
	//}

	// Create TCP client
	opts := mqtt.NewClientOptions()
	//opts.SetTLSConfig(tlsConfig)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MQTT.Broker, cfg.MQTT.Port))
	opts.SetClientID(fmt.Sprintf("%s-sub", cfg.MQTT.ClientId))
	//opts.SetUsername(cfg.MQTT.Username)
	//opts.SetPassword(cfg.MQTT.Password)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Println("Connected established")
	})
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)

	return &Subscriber{
		Client: client,
		Logger: logger,
		Config: cfg,
	}
}

func (s *Subscriber) Start() {
	if token := s.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}
	if token := s.Client.Subscribe(s.Config.MQTT.RawDataTopic, 2, messageSubHandler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	}
	fmt.Println("Subscribed to data topic:", s.Config.MQTT.RawDataTopic)
	if token := s.Client.Subscribe(s.Config.MQTT.MessageTopic, 2, messageSubHandler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	}
	fmt.Println("Subscribed to message topic:", s.Config.MQTT.MessageTopic)

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Client.Unsubscribe(s.Config.MQTT.RawDataTopic)
	s.Client.Unsubscribe(s.Config.MQTT.MessageTopic)
	s.Client.Disconnect(250)
}
