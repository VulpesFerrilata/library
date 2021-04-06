package config

type ServiceSettings struct {
	TranslationFolderPath string   `mapstructure:"TranslationFolderPath"`
	AllowedOrigins        []string `mapstructure:"AllowedOrigins"`
}
