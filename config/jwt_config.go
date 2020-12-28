package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const jwtConfigFilename = "jwt-config"

type JwtConfig struct {
	AccessTokenSettings  TokenSettings `mapstructure:"AccessTokenSettings"`
	RefreshTokenSettings TokenSettings `mapstructure:"RefreshTokenSettings"`
}

func NewJwtConfig() (*JwtConfig, error) {
	jwtCfg := new(JwtConfig)

	v := viper.New()
	v.SetConfigName(jwtConfigFilename)
	v.AddConfigPath("./")
	v.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config.NewJwtConfig")
	}

	if err := v.Unmarshal(jwtCfg); err != nil {
		return nil, errors.Wrap(err, "config.NewJwtConfig")
	}

	return jwtCfg, nil
}
