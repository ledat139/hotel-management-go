package admin_usecase

import (
	"context"
	"hotel-management/internal/dto"
	"hotel-management/internal/repository"
)

type StatUseCase struct {
	statRepo repository.StatRepository
}

func NewStatUseCase(statRepo repository.StatRepository) *StatUseCase {
	return &StatUseCase{statRepo: statRepo}
}

func (u *StatUseCase) GetDashboardStatistics(ctx context.Context) (*dto.StatisticDashboard, error) {
	return u.statRepo.GetDashboardStatistics(ctx)
}
