package requests

type UpdateMe struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=55"`
	LastName  string `json:"last_name" validate:"required,min=3,max=55"`
}
