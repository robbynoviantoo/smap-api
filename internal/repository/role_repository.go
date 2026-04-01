package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM roles ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) CreateRole(ctx context.Context, role *model.Role) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO roles (name) VALUES (?)", role.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RoleRepository) UpdateRole(ctx context.Context, role *model.Role) error {
	res, err := r.db.ExecContext(ctx, "UPDATE roles SET name = ? WHERE id = ?", role.Name, role.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RoleRepository) DeleteRole(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}