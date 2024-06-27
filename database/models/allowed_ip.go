package models

type AllowedIP struct {
	Base
	HasUser
	IP           string
	Active       bool
	Verification string
}
