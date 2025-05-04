package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Gateway struct {
		Host string
		Port string
	}
	MQTT struct {
		Broker   string
		Port     string
		ClientID string
		Username string
		Password string
		Topic    string
		CertPath string
	}
	InfluxDB struct {
		URL    string
		Token  string
		Org    string
		Bucket string
	}
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Unable to decode config: %s", err)
		return nil, err
	}

	return &config, nil
}
