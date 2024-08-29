package models

import (
	"github.com/wizzldev/chat/pkg/encryption"
	"gorm.io/gorm"
	"log"
)

type Message struct {
	Base
	HasGroup
	HasMessageSender
	HasMessageReply
	HasMessageLikes
	Content   string `json:"content"`
	Type      string `json:"type"`
	DataJSON  string `json:"data_json"`
	Encrypted bool   `json:"-"`
}

func (m *Message) AfterFind(*gorm.DB) error {
	if !m.Encrypted {
		return nil
	}
	var err error
	m.Content, err = encryption.DecryptMessage(m.Content)
	return err
}

func (m *Message) BeforeCreate(*gorm.DB) error {
	content, err := encryption.EncryptMessage(m.Content)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	m.Encrypted = true
	m.Content = content
	return nil
}
