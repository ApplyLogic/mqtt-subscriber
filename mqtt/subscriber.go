package mqtt1

import (
	"flag"
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	"github.com/ApplyLogic/mqtt-subscriber/internal/middleware"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Subscriber struct {
	Client mqtt.Client
	Logger *middleware.Logger
	Config *config.Config
}

func New(cfg *config.Config, logger *middleware.Logger) *Subscriber {

	broker := fmt.Sprintf("tls://%s:%d", cfg.MQTT.Broker, cfg.MQTT.Port)
	clientID := cfg.MQTT.ClientId

	//tlsConfig, err := configureTLS(cfg)
	//if err != nil {
	//	log.Fatalf("Error configuring TLS: %v", err)
	//}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	//opts.SetTLSConfig(tlsConfig)

	// Set authentication credentials
	opts.SetUsername(cfg.MQTT.Username)
	opts.SetPassword(cfg.MQTT.Password)

	// Set other options
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Println("Connected established")
	})
	//opts.OnConnect = connectHandler
	//opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)

	return &Subscriber{
		Client: client,
		Logger: logger,
		Config: cfg,
	}
}

func (s *Subscriber) Start() {

	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")
	receiveCount := 0
	choke := make(chan [2]string)

	if token := s.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}

	if token := s.Client.Subscribe(s.Config.MQTT.Topic, 1, subscriber); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for receiveCount < *num {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}

	s.Client.Disconnect(250)
	fmt.Println("Sample Subscriber Disconnected")

	////defer client.Disconnect(250)
	//
	////if token := s.Client.Subscribe(s.Config.MQTT.Topic, 1, messageSubHandler); token.Wait() && token.Error() != nil {
	////	panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	////}
	////fmt.Println("Subscribed to data topic:", s.Config.MQTT.Topic)
	//
	//// Wait for a signal to exit the program gracefully
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//<-sigChan
	////
	////s.Client.Unsubscribe(s.Config.MQTT.Topic)
	////s.Client.Disconnect(250)
}
