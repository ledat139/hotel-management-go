package repository

import (
	"context"
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type StatRepository interface {
	GetDashboardStatistics(ctx context.Context) (*dto.StatisticDashboard, error)
}

type statRepository struct {
	db *gorm.DB
}

func NewStatRepository(db *gorm.DB) StatRepository {
	return &statRepository{db: db}
}
func (r *statRepository) GetDashboardStatistics(ctx context.Context) (*dto.StatisticDashboard, error) {
	var stat dto.StatisticDashboard

	if err := r.db.WithContext(ctx).Model(&models.Room{}).Count(&stat.TotalRooms).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("role = ?", constant.CUSTOMER).
		Count(&stat.TotalCustomers).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Count(&stat.TotalBookings).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Select("SUM(total_price)").
		Where("is_paid = ?", true).
		Scan(&stat.TotalRevenue).Error; err != nil {
		return nil, err
	}

	return &stat, nil
}
