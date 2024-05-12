package requests

type Email struct {
	Email string `json:"email" validate:"required,email"`
}
