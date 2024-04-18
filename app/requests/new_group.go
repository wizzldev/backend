package requests

type NewGroup struct {
	Name    string   `json:"name" validate:"required,min=3,max=55"`
	UserIDs []string `json:"user_ids" validator:""`
}
