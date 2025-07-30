package admin

import (
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	billUseCase *admin_usecase.BillUseCase
}

func NewBillHandler(billUseCase *admin_usecase.BillUseCase) *BillHandler {
	return &BillHandler{billUseCase: billUseCase}
}

func (h *BillHandler) ListBills(c *gin.Context) {
	userName := c.Query("user_name")
	bookingIDStr := c.Query("booking_id")
	exportDate := c.Query("export_date")

	var bookingID int
	if bookingIDStr != "" {
		var err error
		bookingID, err = strconv.Atoi(bookingIDStr)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"Title": "title.bill_management",
				"error": utils.T(c, "error.invalid_room_id"),
				"T":     utils.TmplTranslateFromContext(c),
			})
			return
		}
	}

	bills, err := h.billUseCase.GetFilteredBills(c.Request.Context(), userName, bookingID, exportDate)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "title.bill_management",
			"error": utils.T(c, err.Error()),
			"T":     utils.TmplTranslateFromContext(c),
		})
		return
	}

	c.HTML(http.StatusOK, "bill.html", gin.H{
		"Title": "title.bill_management",
		"Bills": bills,
		"T":     utils.TmplTranslateFromContext(c),
		"Query": gin.H{
			"UserName":   userName,
			"BookingID":  bookingIDStr,
			"ExportDate": exportDate,
		},
	})
}
