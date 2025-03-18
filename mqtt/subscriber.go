package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ApplyLogic/mqtt-subscriber/config"
	"github.com/ApplyLogic/mqtt-subscriber/internal/middleware"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

type Subscriber struct {
	Client mqtt.Client
	Logger *middleware.Logger
}

func New(cfg *config.Config, logger *middleware.Logger) *Subscriber {
	// Load tls cert from your cert file
	cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/test_cert.pem", cfg.TLS.CertPath), fmt.Sprintf("%s/test_cert.key", cfg.TLS.CertPath))
	if err != nil {
		logger.Log("err", err.Error())
		log.Fatal(err)
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
			logger.Log("err", err.Error())
			log.Fatal(err)
		}
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(pemCACert)
		if ok {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			tlsConfig.ClientCAs = certPool
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MQTT.Broker, cfg.MQTT.Port))
	opts.SetClientID(fmt.Sprintf("%s-sub", cfg.MQTT.ClientId))
	//opts.SetTLSConfig(tlsConfig)
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
	}
}

var MessageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}
