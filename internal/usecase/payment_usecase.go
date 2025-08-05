package usecase

import (
	"context"
	"errors"
	"fmt"
	"hotel-management/internal/constant"
	paymentError "hotel-management/internal/error"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"
	"hotel-management/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUseCase struct {
	paymentRepo repository.PaymentRepository
	bookingRepo repository.BookingRepository
	billRepo    repository.BillRepository
}

func NewPaymentUseCase(paymentRepo repository.PaymentRepository, bookingRepo repository.BookingRepository, billRepo repository.BillRepository) *PaymentUseCase {
	return &PaymentUseCase{paymentRepo: paymentRepo, bookingRepo: bookingRepo, billRepo: billRepo}
}

func (u *PaymentUseCase) GetVnPayUrl(ctx context.Context, tx *gorm.DB, bookingID uint, clientIP string) (string, error) {
	paymentURL := ""
	booking, err := u.bookingRepo.GetBookingByIDTx(ctx, tx, bookingID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return paymentURL, paymentError.ErrBookingNotFound
	}
	if err != nil {
		return paymentURL, paymentError.ErrFailedToGetBooking
	}
	if booking.IsPaid {
		return paymentURL, paymentError.ErrBookingHasPaid
	}
	txnRef := fmt.Sprintf("%d-%s", bookingID, uuid.New().String())

	newPayment := &models.Payment{
		BookingID:     booking.ID,
		TransactionID: "",
		PaymentMethod: "vnpay",
		PaymentStatus: constant.PAYMENT_PENDING,
		PaidAt:        time.Now(),
		TxnRef:        txnRef,
	}
	err = u.paymentRepo.CreatePaymentTx(ctx, tx, newPayment)
	if err != nil {
		return paymentURL, errors.New("error.failed_to_save_payment")
	}
	paymentURL, err = utils.CreateVnpayPaymentURL(txnRef, fmt.Sprint(booking.ID), int(booking.TotalPrice), clientIP, constant.HOTEL_ORDER_TYPE)
	if err != nil {
		return paymentURL, errors.New("error.failed_to_create_vnpay_payment")
	}
	return paymentURL, nil
}

func (u *PaymentUseCase) HandleVnpayCallback(ctx context.Context, vnpTxnRef, vnpResponseCode, vnpTransactionNo string) error {
	tx := u.paymentRepo.GetDB()
	return utils.WithTransaction(tx, func(tx *gorm.DB) error {
		payment, err := u.paymentRepo.GetPaymentByTxnRefTx(ctx, tx, vnpTxnRef)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentError.ErrPaymentNotFound
		}
		if err != nil {
			return paymentError.ErrFailedToGetPayment
		}

		booking, err := u.bookingRepo.GetBookingByIDTx(ctx, tx, payment.BookingID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return paymentError.ErrBookingNotFound
		}
		if err != nil {
			return paymentError.ErrFailedToGetBooking
		}

		if payment.PaymentStatus != constant.PAYMENT_PENDING {
			return paymentError.ErrPaymentAlreadyProcessed
		}

		if vnpResponseCode == "00" {
			payment.PaymentStatus = constant.PAYMENT_SUCCESS
			payment.TransactionID = vnpTransactionNo
			booking.IsPaid = true
			booking.BookingStatus = constant.BOOKED

			if err := u.bookingRepo.UpdateBookingTx(ctx, tx, booking); err != nil {
				return paymentError.ErrFailedToUpdateBooking
			}

			bill := &models.Bill{
				BookingID:   booking.ID,
				TotalAmount: booking.TotalPrice,
				ExportAt:    time.Now(),
			}
			if err := u.billRepo.CreateBillTx(ctx, tx, bill); err != nil {
				return paymentError.ErrFailedToCreateBill
			}
		} else {
			payment.PaymentStatus = constant.PAYMENT_FAILED
		}
		if err := u.paymentRepo.UpdatePaymentTx(ctx, tx, payment); err != nil {
			return paymentError.ErrFailedToUpdatePayment
		}
		return nil
	})
}
