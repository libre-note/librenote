package http

type registrationReq struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type updateSettings struct {
	OldPassword     string `json:"old_password" validate:"omitempty,min=8,max=10"`
	NewPassword     string `json:"new_password" validate:"omitempty,min=8,max=100"`
	ListViewEnabled *int8  `json:"list_view_enabled" validate:"required"`
	DarkModeEnabled *int8  `json:"dark_mode_enabled" validate:"required"`
}
