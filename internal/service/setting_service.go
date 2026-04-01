package service

import (
	"context"
	"smap-api/internal/repository"
)

type SettingService struct {
	repo *repository.SettingRepository
}

func NewSettingService(repo *repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetValue(ctx context.Context, key string, defaultValue string) string {
	return s.repo.GetValue(ctx, key, defaultValue)
}
