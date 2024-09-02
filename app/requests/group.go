package requests

import "time"

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

type NewInvite struct {
	MaxUsage   int        `json:"max_usage" validate:"number,min=0,max=50"`
	Expiration *time.Time `json:"expiration" validate:"omitempty,invite_date"`
}

type CustomInvite struct {
	Invite string `json:"invite" validate:"omitempty,min=3,max=15,alphanumunicode"`
}

type Emoji struct {
	Emoji string `json:"emoji" validate:"required,is_emoji"`
}

type Nickname struct {
	Nickname string `json:"nickname" validate:"required,min=3,max=50"`
}
