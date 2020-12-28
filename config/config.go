package config

import (
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "config.NewConfig")
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "config.NewConfig")
	}

	return cfg, nil
}
