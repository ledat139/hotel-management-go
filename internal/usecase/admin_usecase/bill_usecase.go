package admin_usecase

import (
	"context"
	"errors"
	"hotel-management/internal/models"
	"hotel-management/internal/repository"

	"gorm.io/gorm"
)

type BillUseCase struct {
	billRepository repository.BillRepository
}

func NewBillUseCase(billRepository repository.BillRepository) *BillUseCase {
	return &BillUseCase{billRepository: billRepository}
}

func (u *BillUseCase) GetFilteredBills(ctx context.Context, userName string, bookingID int, exportDate string) ([]models.Bill, error) {
	var bills []models.Bill
	bills, err := u.billRepository.SearchBills(ctx, userName, bookingID, exportDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bills, errors.New("error.bill_not_found")
		}
		return bills, errors.New("error.failed_to_get_bill")
	}
	return bills, nil
}
