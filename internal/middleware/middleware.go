package middleware

import (
	"github.com/gin-gonic/gin"
)

const LangKey = "lang"

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = "en"
		}
		c.Set(LangKey, lang)
		c.Next()
	}
}
