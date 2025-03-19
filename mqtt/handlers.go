package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

var messageSubHandler mqtt.MessageHandler = subscriber

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	if token := client.Subscribe("topic/data", 2, messageSubHandler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	}
	fmt.Println("Subscribed to data topic:", "topic/data")
	if token := client.Subscribe("topic/msg", 2, messageSubHandler); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:  %s", token.Error()))
	}
	fmt.Println("Subscribed to message topic:", "topic/msg")
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
