package config

import (
	"awesomeProject4/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-required:"port"`
		BindIP string `yaml:"bind_ip" env-required:"127.0.0.1"`
		Port   string `yaml:"port" env-required:"8080"`
	} `yaml:"listen"`
	MongoDB struct {
		Host       string `json:"host"`
		Port       string `json:"port"`
		Database   string `json:"database"`
		AuthDB     string `json:"auth_Db"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Collection string `json:"collection"`
	} `json:"mongodb"`
}

//Сингл тон
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil) // Выводит что пошло не так
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
