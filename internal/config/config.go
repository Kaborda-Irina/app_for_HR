package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Files properties to work with files
type LoggerConfig struct {
	Level       int    `mapstructure:"level"`
	InfoLogFile string `mapstructure:"info_log_file"`
}

type Mongo struct {
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type CurrencyConfig struct {
	CoefficientEURtoUSD float64 `mapstructure:"coefficientEURtoUSD"`
	CoefficientRUStoUSD float64 `mapstructure:"coefficientRUStoUSD"`
}

type Config struct {
	Port           string `mapstructure:"port"`
	LoggerConfig   `mapstructure:"logger"`
	Mongo          `mapstructure:"mongo"`
	CurrencyConfig `mapstructure:"currency"`
}

func LoadConfig() (config Config, logger *logrus.Logger, err error) {

	//Initialize properties config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	//Initialize logger
	logger = InitLogger(&config.LoggerConfig)

	return config, logger, err
}
