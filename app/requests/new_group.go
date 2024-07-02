package requests

type NewGroup struct {
	Name    string   `json:"name" validate:"required,min=3,max=55"`
	UserIDs []uint   `json:"user_ids" validator:"required,min:3,number"`
	Roles   []string `json:"roles" validate:"required,dive,is_role"`
}

type ModifyRoles struct {
	Roles []string `json:"roles" validate:"required,dive,is_role"`
}

type EditGroupName struct {
	Name string `json:"name" validate:"required,min=3,max=55"`
}
