package service

import (
	"context"
	"smap-api/internal/model"
	"smap-api/internal/repository"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) GetAssets(ctx context.Context, limit, offset int, filter model.AssetFilter) ([]model.Asset, error) {
	return s.repo.GetAssets(ctx, limit, offset, filter)
}

func (s *AssetService) GetAssetsWithCount(ctx context.Context, limit, offset int, filter model.AssetFilter) ([]model.Asset, int, error) {
	assets, err := s.repo.GetAssets(ctx, limit, offset, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountAssets(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}