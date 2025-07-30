package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByTxnRefTx(ctx context.Context, tx *gorm.DB, txnRef string) (*models.Payment, error)
	UpdatePaymentTx(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
	GetDB() *gorm.DB
}
type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}
func (r *paymentRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *paymentRepository) GetPaymentByTxnRefTx(ctx context.Context, tx *gorm.DB, txnRef string) (*models.Payment, error) {
	var payment models.Payment
	err := tx.WithContext(ctx).Where("txn_ref = ?", txnRef).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentTx(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	return tx.WithContext(ctx).Save(payment).Error
}
