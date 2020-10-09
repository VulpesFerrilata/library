package config

import "time"

type TokenSettings struct {
	Alg       string        `mapstructure:"Algorithm"`
	SecretKey string        `mapstructure:"SecretKey"`
	Duration  time.Duration `mapstructure:"Duration"`
}
