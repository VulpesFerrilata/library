package config

type SqlSettings struct {
	DriverName string `mapstructure:"DriverName"`
	DataSource string `mapstructure:"DataSource"`
}
