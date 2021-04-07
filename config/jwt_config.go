package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type JwtConfig struct {
	AccessTokenSettings  TokenSettings `mapstructure:"AccessTokenSettings"`
	RefreshTokenSettings TokenSettings `mapstructure:"RefreshTokenSettings"`
}

func NewJwtConfig(path string) (*JwtConfig, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jwtCfg := new(JwtConfig)

	v := viper.New()

	v.SetConfigFile(absPath)
	viper.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config.NewJwtConfig")
	}

	if err := v.Unmarshal(jwtCfg); err != nil {
		return nil, errors.Wrap(err, "config.NewJwtConfig")
	}

	return jwtCfg, nil
}
