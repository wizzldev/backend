package ws

import (
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/utils"
)

type Map map[string]interface{}

type Message struct {
	Event  string      `json:"event"`
	Data   interface{} `json:"data"`
	HookID string      `json:"hook_id"`
}

type MessageWrapper struct {
	Message  *Message `json:"message"`
	Resource string   `json:"resource"`
}

type ClientMessage struct {
	Content  string `json:"content" validate:"required,min=1,max=500"`
	Type     string `json:"type" validate:"required,max=55"`
	DataJSON string `json:"data_json" validate:"required,json,max=200"`
	HookID   string `json:"hook_id"`
}

type ClientMessageWrapper struct {
	Message  *ClientMessage `json:"message"`
	Resource string         `json:"resource"`
}

func NewMessage(data []byte, conn *Connection) (*ClientMessageWrapper, error) {
	var c ClientMessageWrapper
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	if err := utils.Validator.Struct(c.Message); err != nil {
		go conn.Send(MessageWrapper{
			Message: &Message{
				Event: "error",
				Data:  err.Error(),
			},
			Resource: configs.DefaultWSResource,
		})
		return nil, err
	}
	return &c, nil
}
