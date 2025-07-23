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
}

func NewAdminHandler(authUseCase *admin_usecase.AuthUseCase) *AdminHandler {
	return &AdminHandler{authUseCase: authUseCase}
}

func (h *AdminHandler) AdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Title": "Admin Dashboard",
	})
}

func (h *AdminHandler) AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Admin Login Page",
	})
}

func (h *AdminHandler) HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}

	user, err := h.authUseCase.Login(c.Request.Context(), email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": utils.T(c, err.Error())})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_role", user.Role)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": utils.T(c, "error.internal_server")})
		return
	}

	if user.Role == constant.ADMIN {
		c.Redirect(http.StatusFound, constant.AdminHomePath)
	} else {
		c.Redirect(http.StatusFound, constant.StaffDashboardPath)
	}
}

func (h *AdminHandler) HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, constant.AdminLoginPath)
}
