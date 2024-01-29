package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server Server `mapstructure:"server"`
		Redis  Redis  `mapstructure:"redis"`
		Ninjas Ninjas `mapstructure:"ninjas"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		Timeout  int    `mapstructure:"timeout"`
	}

	Server struct {
		Port int `mapstructure:"port"`
	}

	Ninjas struct {
		ApiKey string `mapstructure:"api-key"`
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

	return &conf, nil
}
