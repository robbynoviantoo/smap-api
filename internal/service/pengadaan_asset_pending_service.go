package service

import (
	"context"
	"errors"
	"smap-api/internal/model"
	"smap-api/internal/repository"
	"time"
)

type PengadaanAssetPendingService struct {
	repo *repository.PengadaanAssetPendingRepository
}

func NewPengadaanAssetPendingService(repo *repository.PengadaanAssetPendingRepository) *PengadaanAssetPendingService {
	return &PengadaanAssetPendingService{repo: repo}
}

func (s *PengadaanAssetPendingService) Create(ctx context.Context, userID uint, req *model.PengadaanAssetCreateRequest) (*model.PengadaanAssetPending, error) {
	if req.Name == "" {
		return nil, errors.New("name wajib diisi")
	}
	if req.Qty <= 0 {
		return nil, errors.New("qty harus lebih dari 0")
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	requestDate := time.Now()
	if req.RequestDate != "" {
		d, err := time.Parse("2006-01-02", req.RequestDate)
		if err != nil {
			return nil, errors.New("format request_date tidak valid: gunakan YYYY-MM-DD")
		}
		requestDate = d
	}

	p := &model.PengadaanAssetPending{
		UserID:      userID,
		RequestDate: requestDate,
		Category:    req.Category,
		Name:        req.Name,
		Merk:        req.Merk,
		Spec:        req.Spec,
		Qty:         req.Qty,
		Unit:        req.Unit,
		Image:       req.Image,
		Remark:      req.Remark,
		Priority:    priority,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PengadaanAssetPendingService) GetAll(ctx context.Context) ([]model.PengadaanAssetPending, error) {
	return s.repo.GetAll(ctx)
}

func (s *PengadaanAssetPendingService) GetByID(ctx context.Context, id uint) (*model.PengadaanAssetPending, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("pengadaan pending not found")
	}
	return p, nil
}

func (s *PengadaanAssetPendingService) Update(ctx context.Context, id uint, req *model.PengadaanAssetUpdateRequest) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("pengadaan pending not found")
	}
	if p.Status != "pending" {
		return errors.New("hanya pengadaan dengan status pending yang dapat diubah")
	}
	return s.repo.Update(ctx, id, req)
}

func (s *PengadaanAssetPendingService) Delete(ctx context.Context, id uint) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("pengadaan pending not found")
	}
	return s.repo.Delete(ctx, id)
}

func (s *PengadaanAssetPendingService) Approve(ctx context.Context, id uint, reviewerID uint) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("pengadaan pending not found")
	}
	if p.Status != "pending" {
		return errors.New("request ini sudah diproses sebelumnya")
	}
	return s.repo.ApproveTx(ctx, p, reviewerID)
}

func (s *PengadaanAssetPendingService) Reject(ctx context.Context, id uint, reviewerID uint, rejectReason string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("pengadaan pending not found")
	}
	if p.Status != "pending" {
		return errors.New("request ini sudah diproses sebelumnya")
	}
	if rejectReason == "" {
		return errors.New("reject_reason wajib diisi saat menolak pengadaan")
	}
	return s.repo.Reject(ctx, id, reviewerID, rejectReason)
}
