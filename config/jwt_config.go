package config

import "github.com/spf13/viper"

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
		return jwtCfg, err
	}

	return jwtCfg, v.Unmarshal(jwtCfg)
}
