package admin

import (
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"hotel-management/internal/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	staffUseCase *admin_usecase.StaffUseCase
}

func NewStaffHandler(staffUseCase *admin_usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{staffUseCase: staffUseCase}
}

func (h *StaffHandler) ListStaffs(c *gin.Context) {
	staffs, err := h.staffUseCase.GetAllStaffs(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.failed_to_get_staff_list"),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.HTML(http.StatusOK, "staff.html", gin.H{
		"Staffs": staffs,
		"T":      utils.TmplTranslateFromContext(c),
		"Title":  "title.staff_management",
	})
}
func (h *StaffHandler) ListCustomers(c *gin.Context) {
	customers, err := h.staffUseCase.GetAllCustomers(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.failed_to_get_customer_list"),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.HTML(http.StatusOK, "customer.html", gin.H{
		"Customers": customers,
		"T":         utils.TmplTranslateFromContext(c),
		"Title":     "title.customer_management",
	})
}

func (h *StaffHandler) CreateStaffPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_staff.html", gin.H{
		"Title": "title.create_staff",
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var form dto.CreateStaffRequest
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "create_staff.html", gin.H{
			"error": utils.T(c, "error.invalid_request"),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}
	err := validator.ValidateCreateStaffInput(form.FullName, form.Phone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_staff.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}
	err = h.staffUseCase.CreateStaff(c, &form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "create_staff.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.StaffManagementPath)
}

func (h *StaffHandler) EditStaffPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_staff_id"),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	staff, err := h.staffUseCase.GetStaffByID(c, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.HTML(http.StatusOK, "edit_staff.html", gin.H{
		"Staff": staff,
		"Title": "title.edit_staff",
		"T":     utils.TmplTranslateFromContext(c),
	})
}

func (h *StaffHandler) EditStaff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_staff_id"),
			"Title": "title.create_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}
	var form dto.UpdateStaffRequest
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "edit_staff.html", gin.H{
			"error": utils.T(c, "error.invalid_request"),
			"Title": "title.edit_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	err = validator.ValidateCreateStaffInput(form.FullName, form.Phone)
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_staff.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.edit_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	err = h.staffUseCase.UpdateStaff(c, &form, id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_staff.html", gin.H{
			"error": utils.T(c, err.Error()),
			"Title": "title.edit_staff",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.StaffManagementPath)
}

func (h *StaffHandler) DeleteStaff(c *gin.Context) {
	staffIDStr := c.Param("id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": utils.T(c, "error.invalid_staff_id"),
			"Title": "title.staff_management",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	err = h.staffUseCase.DeleteStaff(c, staffID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": utils.T(c, "error.failed_to_delete_staff"),
			"Title": "title.staff_management",
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.StaffManagementPath)
}
