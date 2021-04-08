package config

type ServiceSettings struct {
	TranslationFolderPath string   `mapstructure:"TranslationFolderPath"`
	Port                  string   `mapstructure:"Port"`
	AllowedOrigins        []string `mapstructure:"AllowedOrigins"`
}
