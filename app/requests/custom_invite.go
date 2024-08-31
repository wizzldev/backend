package requests

type CustomInvite struct {
	Invite string `json:"invite" validate:"omitempty,min=3,max=15,alphanumunicode"`
}
