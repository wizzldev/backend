package requests

type NewPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type SetNewPassword struct {
	Password string `json:"password" validate:"required"`
}
