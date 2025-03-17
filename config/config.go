package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MQTT *MQTT
	TLS  *TLS
}

type MQTT struct {
	ClientId string
	Broker   string
	Port     int
	Topic    string
	Username string
	Password string
}

type TLS struct {
	CertPath   string
	CACertFile string
}

func LoanConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		MQTT: &MQTT{
			ClientId: viper.GetString("MQTT_CLIENT_ID"),
			Broker:   viper.GetString("MQTT_BROKER"),
			Port:     viper.GetInt("MQTT_PORT"),
			Topic:    viper.GetString("MQTT_TOPIC"),
			Username: viper.GetString("MQTT_USERNAME"),
			Password: viper.GetString("MQTT_PASSWORD"),
		},
		TLS: &TLS{
			CertPath: viper.GetString("TLS_CERT_PATH"),
			//if you want clients to authenticate only with certs issued by your CA
			CACertFile: viper.GetString("TLS_CA_CERT_FILE"),
		},
	}
	return cfg, nil
}
