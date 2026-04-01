package service

import (
	"context"
	"errors"
	"smap-api/internal/model"
	"smap-api/internal/repository"
)

type AssetPendingService struct {
	repo      *repository.AssetPendingRepository
	assetRepo *repository.AssetRepository
}

func NewAssetPendingService(repo *repository.AssetPendingRepository, assetRepo *repository.AssetRepository) *AssetPendingService {
	return &AssetPendingService{repo: repo, assetRepo: assetRepo}
}

func (s *AssetPendingService) CreatePending(ctx context.Context, pending *model.AssetPending) error {
	return s.repo.CreatePending(ctx, pending)
}

func (s *AssetPendingService) GetAllPendings(ctx context.Context) ([]model.AssetPending, error) {
	return s.repo.GetAllPendings(ctx)
}

func (s *AssetPendingService) ReviewPendingAsset(ctx context.Context, id uint, req *model.AssetPendingReviewRequest, reviewerID uint) error {
	pending, err := s.repo.GetPendingByID(ctx, id)
	if err != nil {
		return err
	}
	if pending == nil {
		return errors.New("asset pending not found")
	}

	if req.Status != "approved" && req.Status != "rejected" {
		return errors.New("invalid status: must be approved or rejected")
	}

	if req.Status == "approved" {
		// Insert to assets
		assetCreate := &model.AssetCreateRequest{
			Name:            pending.Name,
			Image:           pending.Image,
			AssetCode:       pending.AssetCode,
			NoAsset:         pending.NoAsset,
			Location:        pending.Location,
			Building:        pending.Building,
			Category:        pending.Category,
			SubCategory:     pending.SubCategory,
			Merk:            pending.Merk,
			Size:            pending.Size,
			Unit:            pending.Unit,
			Status:          pending.Status, // using frontend provided status
			AvailableStatus: "available",
			LastMaintenance: pending.LastMaintenance,
			NextMaintenance: pending.NextMaintenance,
			Remarks:         pending.Remarks,
		}

		err = s.assetRepo.CreateAsset(ctx, assetCreate)
		if err != nil {
			return err
		}
	}

	// Update the pending table
	return s.repo.UpdatePendingStatus(ctx, id, req.Status, reviewerID, req.RejectReason)
}
