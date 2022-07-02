package model

type Label struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	UserID    int32  `json:"user_id"`
	IsTrashed int8   `json:"is_trashed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
