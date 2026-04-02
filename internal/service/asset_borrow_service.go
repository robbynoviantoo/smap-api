package service

import (
	"context"
	"errors"
	"smap-api/internal/model"
	"smap-api/internal/repository"
)

type AssetBorrowService struct {
	repo      *repository.AssetBorrowRepository
	assetRepo *repository.AssetRepository
}

func NewAssetBorrowService(repo *repository.AssetBorrowRepository, assetRepo *repository.AssetRepository) *AssetBorrowService {
	return &AssetBorrowService{repo: repo, assetRepo: assetRepo}
}

// RequestBorrow memvalidasi kondisi asset lalu membuat pengajuan borrow (selalu melalui approval).
func (s *AssetBorrowService) RequestBorrow(ctx context.Context, assetID, userID uint, remark string) error {
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	if asset.AvailableStatus != "available" {
		return errors.New("asset tidak dalam status available")
	}

	if asset.Status != "good" {
		return errors.New("hanya asset dengan status \"good\" yang bisa dipinjam")
	}

	// Cek apakah sudah ada pengajuan pending untuk asset ini
	existing, err := s.repo.GetPendingByAssetID(ctx, assetID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("asset ini sedang dalam proses pengajuan peminjaman")
	}

	return s.repo.CreatePendingBorrow(ctx, &model.AssetBorrowPending{
		AssetID: assetID,
		UserID:  userID,
		Remark:  remark,
	})
}

// ReviewBorrow memproses approve/reject pengajuan borrow.
func (s *AssetBorrowService) ReviewBorrow(ctx context.Context, pendingID uint, req *model.BorrowPendingReviewRequest, reviewerID uint) error {
	if req.Status != "approved" && req.Status != "rejected" {
		return errors.New("invalid status: must be approved or rejected")
	}
	if req.Status == "rejected" && req.RejectReason == "" {
		return errors.New("reject_reason wajib diisi saat menolak pengajuan")
	}
	return s.repo.ReviewBorrowTx(ctx, pendingID, req.Status, req.RejectReason, reviewerID)
}

// ReturnAsset memvalidasi status asset lalu mencatat pengembalian secara transaksional.
func (s *AssetBorrowService) ReturnAsset(ctx context.Context, assetID, userID uint, remark string) error {
	asset, err := s.assetRepo.GetAssetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	if asset.AvailableStatus != "borrowed" {
		return errors.New("asset harus dalam status borrowed untuk di-return")
	}

	return s.repo.ReturnAssetTx(ctx, assetID, userID, remark)
}

// GetAllPendings mengembalikan semua pengajuan borrow.
func (s *AssetBorrowService) GetAllPendings(ctx context.Context) ([]model.AssetBorrowPending, error) {
	return s.repo.GetAllPendingBorrows(ctx)
}
