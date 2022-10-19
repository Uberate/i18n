package config

// I18nConfig is the i18n server config.
type I18nConfig struct {
	WebConfig `json:"web_config" yaml:"web_config" mapstructure:"web_config"`
}
