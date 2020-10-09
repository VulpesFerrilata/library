package config

import (
	"github.com/spf13/viper"
)

const configFilename = "config"

type Config struct {
	ServiceSettings `mapstructure:"ServiceSettings"`
	SqlSettings     `mapstructure:"SqlSettings"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	v := viper.New()
	v.SetConfigName(configFilename)
	v.AddConfigPath("./")
	v.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return cfg, err
	}

	return cfg, v.Unmarshal(cfg)
}
