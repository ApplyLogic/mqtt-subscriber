package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golanshy/tcp-mqtt-service/config"
	"log"
	"time"
)

type Publisher struct {
	client mqtt.Client
	topic  string
}

func NewPublisher(cfg *config.Config) (*Publisher, error) {

	broker := fmt.Sprintf("tls://%s:%s", cfg.MQTT.Broker, cfg.MQTT.Port)
	clientID := cfg.MQTT.ClientID
	topic := cfg.MQTT.Topic

	message := "Hello, secure MQTT!"

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

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}
	//defer client.Disconnect(250)

	token := client.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Error publishing message: %v", token.Error())
	}

	fmt.Println("Message published successfully!")

	return &Publisher{
		client: client,
		topic:  topic,
	}, nil
}

//func configureTLS(cfg *config.Config) (*tls.Config, error) {
//	caCert, err := os.ReadFile("/Users/golanshay/Workspace/tls-certs/ca-cert.pem")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	certFile := fmt.Sprintf("%s/client-cert.pem", cfg.MQTT.CertPath)
//	keyFile := fmt.Sprintf("%s/client-key.pem", cfg.MQTT.CertPath)
//
//	caCertPool := x509.NewCertPool()
//	caCertPool.AppendCertsFromPEM(caCert)
//	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return &tls.Config{
//		RootCAs:      caCertPool,
//		Certificates: []tls.Certificate{cert},
//	}, nil
//}

func (p *Publisher) Publish(message []byte) error {
	token := p.client.Publish(p.topic, 1, false, message)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}
	log.Printf("Published message to topic: %s %s", p.topic, message)
	return nil
}

func (p *Publisher) Close() {
	p.client.Disconnect(250)
}
