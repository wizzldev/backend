package models

type Storage struct {
	Base
	HasUser
	FilePath string `json:"file_path"`
	IsPublic bool   `json:"-"`
}
