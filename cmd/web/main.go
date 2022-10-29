package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/internal/web"
	files2 "github.com/uberate/i18n/pkg/files"
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

	i, err := files2.FromFiles(provider.ISO6391, configInstance.ApplicationConfig.Files...)
	if err != nil {
		panic(err)
	}

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
