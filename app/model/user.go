package model

import "time"

type User struct {
	ID       int32  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	// password hash
	Hash string `json:"hash"`
	// password salt
	Salt            string    `json:"salt"`
	IsActive        int8      `json:"is_active"`
	IsTrashed       int8      `json:"is_trashed"`
	ListViewEnabled int8      `json:"list_view_enabled"`
	DarkModeEnabled int8      `json:"dark_mode_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
