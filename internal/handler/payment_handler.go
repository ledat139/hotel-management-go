package handler

import (
	"errors"
	paymentError "hotel-management/internal/error"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{paymentUseCase: paymentUseCase}
}

func (h *PaymentHandler) GetVnPayUrl(c *gin.Context) {
	bookingIDStr := c.Param("id")
	clientIP := c.ClientIP()

	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil || bookingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_booking_id")})
		return
	}
	if clientIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, "error.invalid_client_ip")})
		return
	}
	paymentURL, err := h.paymentUseCase.GetVnPayUrl(c.Request.Context(), uint(bookingID), clientIP)
	if err != nil {
		switch {
		case errors.Is(err, paymentError.ErrBookingNotFound), errors.Is(err, paymentError.ErrBookingHasPaid):
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrFailedToGetBooking):
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.failed_to_pay_booking")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"payment_url": paymentURL,
	})
}

func (h *PaymentHandler) HandleVnpayCallback(c *gin.Context) {
	vnpTxnRef := c.Query("vnp_TxnRef")
	vnpResponseCode := c.Query("vnp_ResponseCode")
	vnpTransactionNo := c.Query("vnp_TransactionNo")

	if vnpTxnRef == "" || vnpResponseCode == "" || vnpTransactionNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error.invalid_vnpay_callback_parameters"})
		return
	}

	err := h.paymentUseCase.HandleVnpayCallback(c.Request.Context(), vnpTxnRef, vnpResponseCode, vnpTransactionNo)
	if err != nil {
		switch {
		case errors.Is(err, paymentError.ErrPaymentNotFound), errors.Is(err, paymentError.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrFailedToGetPayment), errors.Is(err, paymentError.ErrFailedToGetBooking):
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrPaymentAlreadyProcessed):
			c.JSON(http.StatusBadRequest, gin.H{"error": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrFailedToUpdatePayment), errors.Is(err, paymentError.ErrFailedToUpdateBooking), errors.Is(err, paymentError.ErrFailedToCreateBill):
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.T(c, "error.failed_to_process_payment")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.payment_processed")})
}
