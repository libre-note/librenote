package model

import "time"

type User struct {
	ID       int32  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	// password hash
	Hash string `json:"hash"`
	// password salt
	Salt      string    `json:"salt"`
	IsActive  bool      `json:"is_active"`
	IsTrashed bool      `json:"is_trashed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
