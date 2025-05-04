package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	mqtt2 "github.com/ApplyLogic/mqtt-subscriber/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

func subscriber(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var messageSubHandler mqtt.MessageHandler = mqtt2.subscriber

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("connectHandler: Connected")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	if token := client.Subscribe("tt_controller/device/raw", 1, mqtt2.subscriber); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
	}
	fmt.Println("Subscribed to data topic:", "tt_controller/device/raw")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func configureTLS(cfg *config.Config) (*tls.Config, error) {
	// Load tls cert from your cert file
	cert, err := tls.LoadX509KeyPair(cfg.TLS.ClientCertFile, cfg.TLS.ClientKeyFile)
	if err != nil {
		return nil, err
	}

	// Basic TLS Config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Optionally, if you want clients to authenticate only with certs issued by your CA,
	// you might want to use something like this:
	if cfg.TLS.CACertFile != "" {
		pemCACert, err := os.ReadFile(cfg.TLS.CACertFile)
		if err != nil {
			return nil, err
		}
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(pemCACert)
		if ok {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = certPool
		}
	}
	return tlsConfig, nil
}
