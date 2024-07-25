package requests

type Emoji struct {
	Emoji string `json:"emoji" validate:"required,is_emoji"`
}
