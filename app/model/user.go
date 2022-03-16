package model

import (
	"context"
	"time"
)

type User struct {
	ID       int32  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	// password salted hash
	Hash            string    `json:"hash"`
	IsActive        int8      `json:"is_active"`
	IsTrashed       int8      `json:"is_trashed"`
	ListViewEnabled int8      `json:"list_view_enabled"`
	DarkModeEnabled int8      `json:"dark_mode_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserRepository represent the user's repository contract
type UserRepository interface {
	CreateUser(tx context.Context, user *User) error
	GetUser(tx context.Context, id int32) (User, error)
	GetUserByEmail(tx context.Context, email string) (User, error)
	UpdateUser(tx context.Context, user *User) error
}
