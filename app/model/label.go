package model

import (
	"time"
)

type Label struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	UserID    int32     `json:"user_id"`
	IsTrashed int8      `json:"is_trashed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
