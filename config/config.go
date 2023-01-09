package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Version = "dev"

type Config struct {
	*viper.Viper
}

// GlobalConfig 默认全局配置
var GlobalConfig *Config

// Init 使用 ./application.yaml 初始化全局配置
func Init() {
	GlobalConfig = &Config{
		viper.New(),
	}
	GlobalConfig.AutomaticEnv()
	GlobalConfig.SetConfigName("application")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath(".")
	GlobalConfig.AddConfigPath("./config")

	err := GlobalConfig.ReadInConfig()
	if err != nil {
		logrus.WithField("config", "GlobalConfig").WithError(err).Panicf("unable to read global config")
	}

	GlobalConfig.WatchConfig() // 自动更新配置
	GlobalConfig.OnConfigChange(func(e fsnotify.Event) {
		err := GlobalConfig.ReadInConfig()
		if err == nil {
			logrus.WithField("config", "GlobalConfig").Info("config updated")
		}
	})
}
