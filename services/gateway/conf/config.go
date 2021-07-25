package conf

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/klintcheng/kim/logger"
	"github.com/spf13/viper"
)

// Config Config
type Config struct {
	ServiceID     string   `envconfig:"serviceId"`
	ServiceName   string   `envconfig:"serviceName"`
	Namespace     string   `envconfig:"namespace"`
	Listen        string   `envconfig:"listen"`
	PublicAddress string   `envconfig:"publicAddress"`
	PublicPort    int      `envconfig:"publicPort"`
	Tags          []string `envconfig:"tags"`
	ConsulURL     string   `envconfig:"consulURL"`
}

// Init InitConfig
func Init(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/conf")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}
	logger.Info(config)

	return &config, nil
}
