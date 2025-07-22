package middleware

import (
	"hotel-management/internal/constant"
	"net/http"

	"hotel-management/internal/repository"
	"hotel-management/internal/utils"
	"strings"

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
func RequireAuth(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error.missing_token"})
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error.invalid_token"})
			c.Abort()
			return
		}
		tokenStr := tokenParts[1]

		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.T(c, "error.invalid_token")})
			c.Abort()
			return
		}

		user, err := userRepo.GetUserByEmail(c.Request.Context(), claims.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": utils.T(c, "error.get_user_failed")})
			c.Abort()
			return
		}

		if user.Role != "customer" {
			c.JSON(http.StatusForbidden, gin.H{"error": utils.T(c, "error.access_restricted_to_customers_only")})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{"error": utils.T(c, "error.account_is_not_activated")})
			c.Abort()
			return
		}

		c.Set("userEmail", user.Email)
		c.Set("userName", user.Name)
		c.Next()
	}
}
