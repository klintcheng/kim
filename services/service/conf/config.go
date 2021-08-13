package conf

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kelseyhightower/envconfig"
	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config Config
type Config struct {
	ServiceID     string   `envconfig:"serviceId"`
	NodeID        int64    `envconfig:"nodeId"`
	Listen        string   `envconfig:"listen"`
	PublicAddress string   `envconfig:"publicAddress"`
	PublicPort    int      `envconfig:"publicPort"`
	Tags          []string `envconfig:"tags"`
	ConsulURL     string   `envconfig:"consulURL"`
	RedisAddrs    string   `envconfig:"redisAddrs"`
	BaseDb        string   `envconfig:"baseDb"`
	MessageDb     string   `envconfig:"messageDb"`
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
	if config.PublicAddress == "" {
		config.PublicAddress = kim.GetLocalIP()
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

func MakeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	ac := accesslog.File("./access.log")
	// Remove this line to disable logging to console:
	ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	// ac.SetFormatter(&accesslog.JSON{
	// 	Indent:    "  ",
	// 	HumanTime: true,
	// })
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}
