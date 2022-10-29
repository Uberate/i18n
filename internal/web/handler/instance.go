package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/pkg/provider"
	"net/http"
)

func InstanceGet(config config.I18nConfig, i18n *provider.I18n) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, i18n)
	}
}
