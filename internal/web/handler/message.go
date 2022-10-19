package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uberate/i18n/cmd/web/config"
	"github.com/uberate/i18n/pkg/provider"
	"net/http"
	"strings"
)

func MessageGet(config config.I18nConfig, i18n *provider.I18n) gin.HandlerFunc {
	return func(context *gin.Context) {
		language := context.Param("ln")
		scopes := context.Param("scopes")
		if len(scopes) > 0 && strings.HasSuffix(scopes, "/") {
			scopes = scopes[:len(scopes)-1]
		}
		if len(scopes) > 0 && strings.HasPrefix(scopes, "/") {
			scopes = scopes[1:]
		}

		value, ok := i18n.MessageByString(language, strings.Split(scopes, "/")...)
		if !ok {
			if config.ApplicationConfig.NotFoundWith404 {
				context.JSON(http.StatusNotFound, nil)
				return
			} else {
				value = fmt.Sprintf("%s_%s", language, scopes)
			}
		}
		context.JSON(http.StatusOK, value)
	}
}
