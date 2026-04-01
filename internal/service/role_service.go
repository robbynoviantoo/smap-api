package service

import (
	"context"
	"smap-api/internal/model"
	"smap-api/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.repo.GetAllRoles(ctx)
}

func (s *RoleService) CreateRole(ctx context.Context, role *model.Role) error {
	return s.repo.CreateRole(ctx, role)
}

func (s *RoleService) UpdateRole(ctx context.Context, role *model.Role) error {
	return s.repo.UpdateRole(ctx, role)
}

func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	
	return s.repo.DeleteRole(ctx, id)
}
