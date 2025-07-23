package middleware

import (
	"hotel-management/internal/constant"
	"net/http"

	"github.com/gin-contrib/sessions"
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

func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, constant.AdminLoginPath)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	roleMap := make(map[string]bool)
	for _, r := range allowedRoles {
		roleMap[r] = true
	}

	return func(c *gin.Context) {
		session := sessions.Default(c)
		role, ok := session.Get("user_role").(string)
		if !ok || !roleMap[role] {
			c.Redirect(http.StatusFound, constant.AdminLoginPath)
			c.Abort()
			return
		}
		c.Next()
	}
}
