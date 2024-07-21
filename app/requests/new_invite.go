package requests

import "time"

type NewInvite struct {
	MaxUsage   int        `json:"max_usage" validate:"number,min=0,max=50"`
	Expiration *time.Time `json:"expiration" validate:"omitempty,invite_date"`
}
