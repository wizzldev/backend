package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/pkg/utils"
)

type Map map[string]interface{}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ClientMessage struct {
	Content  string `json:"content" validate:"required,min=1,max=300"`
	Type     string `json:"type" validate:"required,max=55"`
	DataJSON string `json:"data_json" validate:"required,json,max=55"`
}

func (c *ClientMessage) Make(data []byte, conn *Connection) error {
	err := json.NewDecoder(bytes.NewReader(data)).Decode(c)
	if err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	if err := utils.Validator.Struct(c); err != nil {
		go conn.Send(Message{
			Event: "error",
			Data:  err.Error(),
		})
		return err
	}
	return nil
}
