package requests

type NewBot struct {
	Name string `json:"name" validate:"required,min=3,max=55"`
}

type ApplicationInvite struct {
	GroupID uint `json:"group_id"`
	BotID   uint `json:"bot_id"`
}
