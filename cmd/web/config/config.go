package config

import "github.com/uberate/mocker-utils/gins"

// I18nConfig is the i18n server config.
type I18nConfig struct {
	WebConfig gins.WebConfig `json:"web_config" yaml:"web_config" mapstructure:"web_config"`

	ApplicationConfig ApplicationConfig `json:"application_config" yaml:"application_config" mapstructure:"application_config"`
}

type ApplicationConfig struct {
	Files           []string `json:"files" yaml:"files" mapstructure:"files"`
	Readonly        bool     `json:"readonly" yaml:"readonly" mapstructure:"readonly"`
	NotFoundWith404 bool     `json:"not_found_with_404" yaml:"not_found_with_404" mapstructure:"not_found_with_404"`
}
