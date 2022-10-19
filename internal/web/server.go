package web

import (
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/internal/web/handler"
	"github.com/uberate/i18n/pkg/provider"
)

var currentVersion = "v1"

func RegisterHandler(engine *gin.Engine, config config.I18nConfig, i18nInstance *provider.I18n) {
	v1 := engine.Group(currentVersion)
	{
		message := v1.Group("message")
		{
			message.GET("ln/:ln/*scopes", handler.MessageGet(config, i18nInstance))
		}

		languages := v1.Group("language")
		v1.GET("languages", handler.LanguageList(config, i18nInstance))
		{
			languages.GET("/standards", handler.StandardList(config, i18nInstance))
			languages.GET("/:language", handler.LanguageGet(config, i18nInstance))
		}
	}
}
