package handler

import (
	"errors"
	paymentError "hotel-management/internal/error"
	"hotel-management/internal/usecase"
	"hotel-management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{paymentUseCase: paymentUseCase}
}

// HandleVnpayCallback godoc
// @Summary      Process VnPay callback
// @Description  Handle the response returned from VnPay after a payment attempt.
// @Tags         payments
// @Param        vnp_TxnRef          query  string  true  "Transaction Reference"
// @Param        vnp_ResponseCode    query  string  true  "VnPay Response Code"
// @Param        vnp_TransactionNo   query  string  true  "VnPay Transaction Number"
// @Success      200  {object}  map[string]string  "Payment processed successfully"
// @Failure      400  {object}  map[string]string  "Invalid callback parameters"
// @Failure      404  {object}  map[string]string  "Payment not found"
// @Failure      409  {object}  map[string]string  "Payment has already been processed or booking already paid"
// @Failure      500  {object}  map[string]string  "Failed to update payment, process payment, or create bill"
// @Router       /payments/vnpay_return [get]
func (h *PaymentHandler) HandleVnpayCallback(c *gin.Context) {
	vnpTxnRef := c.Query("vnp_TxnRef")
	vnpResponseCode := c.Query("vnp_ResponseCode")
	vnpTransactionNo := c.Query("vnp_TransactionNo")

	if vnpTxnRef == "" || vnpResponseCode == "" || vnpTransactionNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error.invalid_vnpay_callback_parameters"})
		return
	}

	err := h.paymentUseCase.HandleVnpayCallback(c.Request.Context(), vnpTxnRef, vnpResponseCode, vnpTransactionNo)
	if err != nil {
		switch {
		case errors.Is(err, paymentError.ErrPaymentNotFound), errors.Is(err, paymentError.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrFailedToGetPayment), errors.Is(err, paymentError.ErrFailedToGetBooking):
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrPaymentAlreadyProcessed):
			c.JSON(http.StatusBadRequest, gin.H{"message": utils.T(c, err.Error())})
		case errors.Is(err, paymentError.ErrFailedToUpdatePayment), errors.Is(err, paymentError.ErrFailedToUpdateBooking), errors.Is(err, paymentError.ErrFailedToCreateBill):
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, err.Error())})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": utils.T(c, "error.failed_to_process_payment")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.T(c, "success.payment_processed")})
}
