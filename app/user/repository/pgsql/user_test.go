package pgsql_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"librenote/app/model"
	userRepo "librenote/app/user/repository/pgsql"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	u := &model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      "2o403w24o32043204weorjwe",
		IsActive:  1,
		UpdatedAt: time.Now().UTC(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "INSERT INTO users"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.FullName, u.Email, u.Hash, u.IsActive, u.UpdatedAt).WillReturnResult(sqlmock.NewResult(11, 1))

	ur := userRepo.NewPgsqlUserRepository(db)
	err = ur.CreateUser(context.TODO(), u)
	assert.NoError(t, err)
	assert.Equal(t, int32(11), u.ID)

}

func TestGetUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUser := model.User{
		ID: 1, FullName: "Mr. Test", Email: "mrtest@example.com", Hash: "2o403w24o32043204weorjwe",
		IsActive: 1, IsTrashed: 0, ListViewEnabled: 1, DarkModeEnabled: 0,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "full_name", "email", "hash", "is_active", "is_trashed", "list_view_enabled", "dark_mode_enabled",
		"created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.FullName, mockUser.Email, mockUser.Hash,
			mockUser.IsActive, mockUser.IsTrashed, mockUser.ListViewEnabled, mockUser.DarkModeEnabled,
			mockUser.CreatedAt, mockUser.UpdatedAt)

	query := "SELECT id, full_name, email, hash, is_active, is_trashed, list_view_enabled, dark_mode_enabled, created_at, updated_at FROM users WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	ur := userRepo.NewPgsqlUserRepository(db)

	num := int32(1)
	user, err := ur.GetUser(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, num, user.ID)
	assert.Equal(t, "mrtest@example.com", user.Email)

}

func TestGetUserByEmail(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "full_name", "email", "hash", "is_active", "is_trashed", "list_view_enabled", "dark_mode_enabled",
		"created_at", "updated_at"}).
		AddRow(1, "Mr. Test", "mrtest@example.com", "skflrrweoiruowiu43", 1, 0, 1, 1, time.Now().UTC(), time.Now().UTC())

	query := "SELECT id, full_name, email, hash, is_active, is_trashed, list_view_enabled, dark_mode_enabled, created_at, updated_at FROM users WHERE email = \\$1 LIMIT 1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	ur := userRepo.NewPgsqlUserRepository(db)

	email := "mrtest@example.com"
	user, err := ur.GetUserByEmail(context.TODO(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Mr. Test", user.FullName)

}

func TestUpdateUser(t *testing.T) {
	u := &model.User{
		ID:              12,
		Hash:            "2o403w24o32043204weorjwe",
		IsActive:        1,
		IsTrashed:       0,
		ListViewEnabled: 0,
		DarkModeEnabled: 0,
		UpdatedAt:       time.Now().UTC(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := "UPDATE users SET hash = \\$2, is_active = \\$3, is_trashed = \\$4, list_view_enabled = \\$5, dark_mode_enabled = \\$6, updated_at = \\$7 WHERE id = \\$1"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(u.ID, u.Hash, u.IsActive, u.IsTrashed, u.ListViewEnabled, u.DarkModeEnabled, u.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(12, 1))

	ur := userRepo.NewPgsqlUserRepository(db)

	err = ur.UpdateUser(context.TODO(), u)
	assert.NoError(t, err)
}
