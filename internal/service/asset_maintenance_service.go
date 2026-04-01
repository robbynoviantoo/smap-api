package service

import (
	"context"
	"errors"
	"smap-api/internal/model"
	"smap-api/internal/repository"
	"time"
)

type AssetMaintenanceService struct {
	repo       *repository.AssetMaintenanceRepository
	assetRepo  *repository.AssetRepository
	settingSvc *SettingService
}

func NewAssetMaintenanceService(repo *repository.AssetMaintenanceRepository, assetRepo *repository.AssetRepository, settingSvc *SettingService) *AssetMaintenanceService {
	return &AssetMaintenanceService{
		repo:       repo,
		assetRepo:  assetRepo,
		settingSvc: settingSvc,
	}
}

func (s *AssetMaintenanceService) ScheduleMaintenance(ctx context.Context, assetID uint, req *model.MaintenanceScheduleRequest, userID uint) error {
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	if asset.AvailableStatus != "available" && asset.AvailableStatus != "Tersedia" {
		return errors.New("Asset tidak bisa diatur jadwal perawatannya karena sedang dipinjam atau diambil")
	}

	// Parse date
	nextDate, err := time.Parse("2006-01-02", req.NextMaintenanceDate)
	if err != nil {
		return errors.New("invalid date format: use YYYY-MM-DD")
	}

	// Validate date >= today (ignoring time by truncating both)
	today := time.Now().Truncate(24 * time.Hour)
	if nextDate.Truncate(24*time.Hour).Before(today) {
		return errors.New("next_maintenance_date must be after or equal to today")
	}

	return s.repo.ScheduleMaintenanceTx(ctx, assetID, nextDate, userID)
}

func (s *AssetMaintenanceService) StartMaintenance(ctx context.Context, assetID uint, req *model.MaintenanceStartRequest, userID uint) (bool, error) {
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return false, err
	}
	if asset == nil {
		return false, errors.New("asset not found")
	}

	if asset.AvailableStatus != "available" && asset.AvailableStatus != "Tersedia" {
		return false, errors.New("Asset tidak bisa dimaintenance karena sedang dipinjam atau diambil")
	}

	if asset.Status != "good" && asset.Status != "Baru" && asset.Status != "Baik" {
		// Loosening the strict "good" since the user might use "Baru" "Tersedia"
		// adjust condition based on actual usage, assumed "maintenance" means it's already in maintenance
		if asset.Status == "maintenance" || asset.Status == "Maintenance" {
			return false, errors.New("Asset sudah dalam status maintenance")
		}
	}

	// Check existing pending request
	existing, err := s.repo.GetPendingStartByAssetID(ctx, assetID)
	if err != nil {
		return false, err
	}
	if existing != nil {
		return false, errors.New("Asset ini sudah memiliki pengajuan maintenance yang menunggu approval")
	}

	approvalEnabled := s.settingSvc.GetValue(ctx, "asset_maintenance_approval_enabled", "true")

	if approvalEnabled == "true" {
		pendingReq := &model.AssetMaintenancePending{
			AssetID: assetID,
			UserID:  userID,
			Notes:   req.Notes,
		}
		err = s.repo.CreatePendingRequest(ctx, pendingReq)
		return true, err // true indicates it needs approval
	}

	// Approval disabled = Start directly
	err = s.repo.StartMaintenanceTx(ctx, assetID, userID, req.Notes)
	return false, err
}

func (s *AssetMaintenanceService) ReviewPendingStart(ctx context.Context, pendingID uint, req *model.MaintenancePendingReviewRequest, reviewerID uint) error {
	if req.Status != "approved" && req.Status != "rejected" {
		return errors.New("invalid status: must be approved or rejected")
	}

	// The logic for approval -> starting maintenance is handled in the tx
	return s.repo.ReviewPendingTx(ctx, pendingID, req.Status, req.RejectReason, reviewerID)
}

func (s *AssetMaintenanceService) FinishMaintenance(ctx context.Context, assetID uint, req *model.MaintenanceFinishRequest, userID uint) error {
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	if asset.Status != "maintenance" && asset.Status != "Maintenance" {
		return errors.New("Asset sedang tidak dalam status maintenance")
	}

	var parsedDate *time.Time
	if req.NextMaintenanceDate != nil && *req.NextMaintenanceDate != "" {
		d, err := time.Parse("2006-01-02", *req.NextMaintenanceDate)
		if err != nil {
			return errors.New("invalid date format: use YYYY-MM-DD")
		}
		today := time.Now().Truncate(24 * time.Hour)
		if d.Truncate(24*time.Hour).Before(today) {
			return errors.New("next_maintenance_date must be after or equal to today")
		}
		parsedDate = &d
	}

	return s.repo.FinishMaintenanceTx(ctx, assetID, asset.NextMaintenance, parsedDate, userID, req.Notes)
}

func (s *AssetMaintenanceService) GetAllPendingRequests(ctx context.Context) ([]model.AssetMaintenancePending, error) {
	return s.repo.GetAllPendingRequests(ctx)
}
