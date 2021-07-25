package conf

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/kelseyhightower/envconfig"
	"github.com/klintcheng/kim/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
}

// Config Config
type Config struct {
	ServiceID     string   `envconfig:"serviceId"`
	Namespace     string   `envconfig:"namespace"`
	Listen        string   `envconfig:"listen"`
	PublicAddress string   `envconfig:"publicAddress"`
	PublicPort    int      `envconfig:"publicPort"`
	Tags          []string `envconfig:"tags"`
	ConsulURL     string   `envconfig:"consulURL"`
	RedisAddrs    string   `envconfig:"redisAddrs"`
	RpcURL        string   `envconfig:"ppcURL"`
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

func InitRedis(addr string, pass string) (*redis.Client, error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return redisdb, nil
}

// InitFailoverRedis init redis with sentinels
func InitFailoverRedis(masterName string, sentinelAddrs []string, password string, timeout time.Duration) (*redis.Client, error) {
	redisdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: sentinelAddrs,
		Password:      password,
		DialTimeout:   time.Second * 5,
		ReadTimeout:   timeout,
		WriteTimeout:  timeout,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		logrus.Warn(err)
	}
	return redisdb, nil
}
