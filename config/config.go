package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MQTT *MQTT
	TLS  *TLS
}

type MQTT struct {
	ClientId     string
	Broker       string
	Port         int
	RawDataTopic string
	MessageTopic string
	Username     string
	Password     string
}

type TLS struct {
	CACertFile     string
	ClientCertFile string
	ClientKeyFile  string
}

func LoanConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		MQTT: &MQTT{
			ClientId:     viper.GetString("MQTT_CLIENT_ID"),
			Broker:       viper.GetString("MQTT_BROKER"),
			Port:         viper.GetInt("MQTT_PORT"),
			RawDataTopic: viper.GetString("MQTT_RAW_TOPIC"),
			MessageTopic: viper.GetString("MQTT_MSG_TOPIC"),
			Username:     viper.GetString("MQTT_USERNAME"),
			Password:     viper.GetString("MQTT_PASSWORD"),
		},
		TLS: &TLS{
			ClientCertFile: viper.GetString("TLS_CLIENT_CERT_FILE"),
			ClientKeyFile:  viper.GetString("TLS_CLIENT_KEY_FILE"),
			//if you want clients to authenticate only with certs issued by your CA
			CACertFile: viper.GetString("TLS_CA_CERT_FILE"),
		},
	}
	return cfg, nil
}
