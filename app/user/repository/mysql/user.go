package mysql

import (
	"context"
	"database/sql"
	"errors"
	"librenote/app/model"
)

type userRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) model.UserRepository {
	return &userRepository{
		db: db,
	}
}

const createUser = `INSERT INTO users (
  full_name, email, hash, is_active, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?
)
`

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	stmt, err := r.db.PrepareContext(ctx, createUser)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		user.FullName,
		user.Email,
		user.Hash,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

const getUser = `SELECT id, full_name, email, hash, is_active, is_trashed, list_view_enabled, dark_mode_enabled,
created_at, updated_at FROM users WHERE id = ? LIMIT 1
`

func (r *userRepository) GetUser(ctx context.Context, id int32) (model.User, error) {
	row := r.db.QueryRowContext(ctx, getUser, id)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Hash,
		&i.IsActive,
		&i.IsTrashed,
		&i.ListViewEnabled,
		&i.DarkModeEnabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `SELECT id, full_name, email, hash, is_active, is_trashed, list_view_enabled, dark_mode_enabled,
created_at, updated_at FROM users WHERE email = ? LIMIT 1
`

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	row := r.db.QueryRowContext(ctx, getUserByEmail, email)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Hash,
		&i.IsActive,
		&i.IsTrashed,
		&i.ListViewEnabled,
		&i.DarkModeEnabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `UPDATE users
SET hash = ?,
is_active = ?,
is_trashed = ?,
list_view_enabled = ?,
dark_mode_enabled = ?,
updated_at = ?
WHERE id = ?
`

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	stmt, err := r.db.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx,
		user.ID,
		user.Hash,
		user.IsActive,
		user.IsTrashed,
		user.ListViewEnabled,
		user.DarkModeEnabled,
		user.UpdatedAt,
	)

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		return errors.New("nothing changed")
	}

	return nil
}
