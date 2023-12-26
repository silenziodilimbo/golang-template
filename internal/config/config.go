package config

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ListenConfig struct {
	IP   string `yaml:"IP"`
	Port int    `yaml:"Port"`
}

type GSConfig struct {
	GSEvents string `yaml:"GSEvents"`
}

type AgentConfig struct {
	AgentEvents string `yaml:"AgentEvents"`
}

type TopicConfig struct {
	GS    GSConfig    `yaml:"GS"`
	Agent AgentConfig `yaml:"Agent"`
}

type Config struct {
	Listen ListenConfig `yaml:"Listen"`
	Topic  TopicConfig  `yaml:"Topic"`
}

var once sync.Once
var config *Config
var err error

func GetConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}

		viperConfig := viper.New()

		viperConfig.AddConfigPath("./")
		viperConfig.SetConfigName("config")
		viperConfig.SetConfigType("yaml") //设置文件的类型

		err = viperConfig.ReadInConfig()
		if err != nil {
			logrus.Errorf("read config failed: %v", err)
			return
		}
		err = viperConfig.Unmarshal(&config)
		if err != nil {
			return
		}
		logrus.Infof("config=%v", config)
	})
	return config, err
}
