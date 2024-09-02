package requests

type PushToken struct {
	Token string `json:"token" validate:"required"`
}
