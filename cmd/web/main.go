package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/internal/web"
	"github.com/uberate/i18n/pkg/provider"
	"github.com/uberate/mocker-utils/files"
	"github.com/uberate/mocker-utils/gins"
)

var configInstance = &config.I18nConfig{
	WebConfig: gins.WebConfig{
		Addr: []string{":3000"},
		Mod:  gin.TestMode,
	},
	ApplicationConfig: config.ApplicationConfig{},
}

func main() {
	engine := gin.Default()

	i := provider.NewI18n(provider.Custom)

	web.RegisterHandler(engine, *configInstance, i)

	if err := gins.GinStart(engine, configInstance.WebConfig); err != nil {
		panic(err)
	}
}

func init() {
	if err := files.ReadConfig("",
		"",
		"",
		true,
		"",
		configInstance); err != nil {
		panic(err)
	}
}
