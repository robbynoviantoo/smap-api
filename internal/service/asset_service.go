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

func (s *AssetService) GetAssetByID(ctx context.Context, id uint) (*model.Asset, error) {
	return s.repo.GetAssetByID(ctx, id)
}

func (s *AssetService) CreateAsset(ctx context.Context, asset *model.AssetCreateRequest) error {
	return s.repo.CreateAsset(ctx, asset)
}

func (s *AssetService) UpdateAsset(ctx context.Context, asset *model.AssetUpdateRequest) error {
	return s.repo.UpdateAsset(ctx, asset)
}

func (s *AssetService) DeleteAsset(ctx context.Context, id uint, deletedBy *uint) error {
	return s.repo.DeleteWithBackup(ctx, id, deletedBy)
}

// GetDeletedAssets mengembalikan list asset yang sudah dihapus beserta total.
func (s *AssetService) GetDeletedAssets(ctx context.Context, limit, offset int) ([]model.AssetDelete, int, error) {
	list, err := s.repo.GetDeletedAssets(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountDeletedAssets(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}