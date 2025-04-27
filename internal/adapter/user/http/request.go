package http

type SignUpRequest struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
}

type SignInRequest struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}
