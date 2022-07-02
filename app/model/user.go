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

type UserDetails struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	ListViewEnabled int8   `json:"list_view_enabled"`
	DarkModeEnabled int8   `json:"dark_mode_enabled"`
}

// UserRepository represent the user's repository contract
type UserRepository interface {
	CreateUser(tx context.Context, user *User) error
	GetUser(tx context.Context, id int32) (User, error)
	GetUserByEmail(tx context.Context, email string) (User, error)
	UpdateUser(tx context.Context, user *User) error
}

// Password struct
type Password struct {
	OldPassword string
	NewPassword string
	IsChanged   bool
}

// UserUsecase represent the user's usecase contract
type UserUsecase interface {
	Registration(c context.Context, m *User) (err error)
	Login(c context.Context, email, password string) (token string, err error)
	GetUserDetails(c context.Context, id int32) (user *UserDetails, err error)
	Update(c context.Context, m *User, p Password) error
}
