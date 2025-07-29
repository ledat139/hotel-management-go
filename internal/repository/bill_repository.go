package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type BillRepository interface {
	CreateBillTx(ctx context.Context, tx *gorm.DB, bill *models.Bill) error
}
type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

func (r *billRepository) CreateBillTx(ctx context.Context, tx *gorm.DB, bill *models.Bill) error {
	return tx.WithContext(ctx).Create(&bill).Error
}
