package models

type File struct {
	Base
	Path          string  `json:"path"`
	Name          string  `json:"name"`
	AccessToken   *string `json:"access_token"`
	Type          string  `json:"type"`
	Discriminator string  `json:"discriminator"`
}
