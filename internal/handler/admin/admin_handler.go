package admin

import (
	"hotel-management/internal/constant"
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	authUseCase *admin_usecase.AuthUseCase
	statUseCase *admin_usecase.StatUseCase
}

func NewAdminHandler(authUseCase *admin_usecase.AuthUseCase, statUseCase *admin_usecase.StatUseCase) *AdminHandler {
	return &AdminHandler{authUseCase: authUseCase, statUseCase: statUseCase}
}

func (h *AdminHandler) AdminDashboard(c *gin.Context) {
	stat, err := h.statUseCase.GetDashboardStatistics(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{
			"error": utils.T(c, "failed_to_load_statistics"),
			"Title": "title.dashboard",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"stat":  stat,
		"Title": "title.admin_dashboard",
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *AdminHandler) AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"T":     utils.TmplTranslateFromContext(c),
		"Title": "title.admin_login_page",
	})
}

func (h *AdminHandler) HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": utils.T(c, "error.invalid_request"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.admin_login_page",
		})
		return
	}

	user, err := h.authUseCase.Login(c.Request.Context(), email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": utils.T(c, err.Error()),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.admin_login_page",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": utils.T(c, "error.internal_server"),
			"T":     utils.TmplTranslateFromContext(c),
			"Title": "title.admin_login_page",
		})
		return
	}

	c.Redirect(http.StatusFound, constant.AdminHomePath)
}

func (h *AdminHandler) HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, constant.AdminLoginPath)
}
