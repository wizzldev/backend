package requests

type Register struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=55"`
	LastName  string `json:"last_name" validate:"required,min=3,max=55"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=250"`
}
