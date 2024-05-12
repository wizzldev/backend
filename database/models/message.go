package models

type Message struct {
	Base
	HasGroup
	HasMessageSender
	HasMessageReply
	Content  string `json:"content"`
	Type     string `json:"type"`
	DataJSON string `json:"data_json"`
}