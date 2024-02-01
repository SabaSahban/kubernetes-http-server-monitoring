package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server        Server `mapstructure:"server"`
		Redis         Redis  `mapstructure:"redis"`
		CheckInterval int    `mapstructure:"interval"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		Timeout  int    `mapstructure:"timeout"`
	}

	Server struct {
		Port int `mapstructure:"port"`
	}
)

func Load() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return &Config{}, err
	}

	viper.AutomaticEnv()
	conf.Redis.Password = viper.GetString("REDIS_PASSWORD")

	return &conf, nil
}
