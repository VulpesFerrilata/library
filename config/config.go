package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceSettings `mapstructure:"ServiceSettings"`
	SqlSettings     `mapstructure:"SqlSettings"`
}

func NewConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cfg := new(Config)

	v := viper.New()
	v.SetConfigFile(absPath)
	viper.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, errors.WithStack(err)
	}

	return cfg, nil
}
