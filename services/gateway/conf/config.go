package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/spf13/viper"
)

// Config Config
type Config struct {
	ServiceID     string   `envconfig:"serviceId"`
	ServiceName   string   `envconfig:"serviceName"`
	Listen        string   `envconfig:"listen"`
	PublicAddress string   `envconfig:"publicAddress"`
	PublicPort    int      `envconfig:"publicPort"`
	Tags          []string `envconfig:"tags"`
	ConsulURL     string   `envconfig:"consulURL"`
	AppSecret     string   `envconfig:"appSecret"`
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

	localIP := kim.GetLocalIP()
	if os.Getenv("external_ip") != "" {
		localIP = os.Getenv("external_ip")
	}
	var config = Config{
		ServiceID:     fmt.Sprintf("gate_%s", strings.ReplaceAll(localIP, ".", "")),
		ServiceName:   "gateway",
		Listen:        ":8000",
		PublicAddress: localIP,
		PublicPort:    8000,
		ConsulURL:     fmt.Sprintf("%s:8500", localIP),
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn(err)
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			return nil, err
		}
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	logger.Info(config)
	return &config, nil
}
