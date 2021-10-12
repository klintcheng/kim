package conf

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/spf13/viper"
)

// Config Config
type Config struct {
	ServiceID       string
	ServiceName     string `default:"wgateway"`
	Listen          string `default:":8000"`
	PublicAddress   string
	PublicPort      int `default:"8000"`
	Tags            []string
	Meta            map[string]string
	ConsulURL       string
	MonitorPort     int `default:"8001"`
	AppSecret       string
	LogLevel        string `default:"INFO"`
	MessageGPool    int    `default:"10000"`
	ConnectionGPool int    `default:"15000"`
}

func (c Config) String() string {
	bts, _ := json.Marshal(c)
	return string(bts)
}

// Init InitConfig
func Init(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/conf")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		logger.Warn(err)
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			return nil, err
		}
	}
	err := envconfig.Process("kim", &config)
	if err != nil {
		return nil, err
	}
	if config.ServiceID == "" {
		localIP := kim.GetLocalIP()
		config.ServiceID = fmt.Sprintf("gate_%s", strings.ReplaceAll(localIP, ".", ""))
	}
	if config.PublicAddress == "" {
		config.PublicAddress = kim.GetLocalIP()
	}
	logger.Info(config)
	return &config, nil
}
