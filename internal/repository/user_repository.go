package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)


type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	row := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, is_first_login, created_at FROM users WHERE email = ?`, email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.IsFirstLogin, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	user := &model.User{}
	row := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, image, no_handphone, is_first_login, created_at FROM users WHERE id = ?`, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Image, &user.NoHandphone, &user.IsFirstLogin, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	var users []model.UserResponse
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, first_name, COALESCE(last_name, '-') as last_name, email, COALESCE(no_handphone, '-') as no_handphone, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.UserResponse
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.NoHandphone, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT roles.name
		FROM roles
		JOIN user_roles ON user_roles.role_id = roles.id
		WHERE user_roles.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateOnboarding menyimpan data profil pertama kali login dan set is_first_login = false.
func (r *UserRepository) UpdateOnboarding(ctx context.Context, userID uint, firstName, lastName, noHandphone, hashedPassword, image string) error {
	now := time.Now()
	query := `UPDATE users
		SET first_name = ?,
		    last_name  = ?,
		    no_handphone = ?,
		    password   = ?,
		    image      = CASE WHEN ? != '' THEN ? ELSE image END,
		    is_first_login = 0,
		    updated_at = ?
		WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		firstName, lastName, noHandphone, hashedPassword,
		image, image, // CASE args
		now, userID,
	)
	return err
}
