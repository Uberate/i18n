package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/pkg/provider"
	"net/http"
)

func StandardList(config config.I18nConfig, i18n *provider.I18n) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, []string{
			provider.Custom,
			provider.ISO6391,
			provider.ISO6392B,
			provider.ISO6392T,
			provider.ISO6393,
		})
	}
}

func LanguageList(config config.I18nConfig, i18n *provider.I18n) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, provider.Mapper)
	}
}

func LanguageGet(config config.I18nConfig, i18n *provider.I18n) gin.HandlerFunc {
	return func(context *gin.Context) {
		keyName := context.Param("language")
		context.JSON(http.StatusOK, provider.GetLanguageKey(keyName))
	}
}
