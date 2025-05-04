package mqtt

import (
	"fmt"
	"github.com/golanshy/tcp-mqtt-service/config"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageHandler func(topic string, message []byte)

type Subscriber struct {
	client mqtt.Client
	topic  string
}

func NewSubscriber(cfg *config.Config, handler MessageHandler) (*Subscriber, error) {

	broker := fmt.Sprintf("tls://%s:%s", cfg.MQTT.Broker, cfg.MQTT.Port)
	clientID := cfg.MQTT.ClientID
	topic := cfg.MQTT.Topic

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
		log.Println("Connected to MQTT broker")
		if token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			log.Printf("msg.Payload: %s", msg.Payload())
			handler(msg.Topic(), msg.Payload())
		}); token.Wait() && token.Error() != nil {
			log.Printf("Error subscribing to topic %s: %v", topic, token.Error())
		} else {
			log.Printf("Subscribed to topic: %s", topic)
		}
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection to MQTT broker lost: %v", err)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("failed to connect to MQTT broker: %w", token.Error())
		} else {
			log.Printf("Connected to MQTT broker")
		}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}
	//defer client.Disconnect(250)

	return &Subscriber{
		client: client,
		topic:  topic,
	}, nil
}

func (s *Subscriber) Close() {
	if token := s.client.Unsubscribe(s.topic); token.Wait() && token.Error() != nil {
		log.Printf("Error unsubscribing from topic %s: %v", s.topic, token.Error())
	}
	s.client.Disconnect(250)
}
