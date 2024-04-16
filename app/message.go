package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Content  string `json:"content" validate:"required,min=1,max=300"`
	Type     string `json:"type" validate:"required,eq=message|eq=colored_message"`
	DataJSON string `json:"data_json" validate:"required,json,max=55"`
}

var ctx = context.Background()

func MessageActionHandler(s *ws.Server, conn *ws.Connection, userID uint, data []byte) error {
	if configs.Env.Debug {
		fmt.Printf("WS[%v] New message: %s by user ID %v\n", s.ID, string(data), userID)
	}

	user, err := getCachedUser(userID)

	if err != nil {
		go conn.Send(ws.Message{
			Event: "error",
			Data:  err.Error(),
		})
		return err
	}

	gID, err := strconv.Atoi(s.ID)
	if err != nil {
		return err
	}

	var msg Message
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&msg)
	if err != nil {
		return fmt.Errorf("failed to decode body: %w", err)
	}

	if err = utils.Validator.Struct(msg); err != nil {
		go conn.Send(ws.Message{
			Event: "error",
			Data:  err.Error(),
		})
		return err
	}

	switch msg.Type {
	case "message":
		message := models.Message{
			HasGroup: models.HasGroup{
				GroupID: uint(gID),
			},
			HasMessageSender: models.HasMessageSender{
				SenderID: userID,
			},
			Content:  msg.Content,
			Type:     msg.Type,
			DataJSON: msg.DataJSON,
		}
		database.DB.Create(&message)

		events.DispatchMessage(s.ID, getCachedGroupUserIDs(s.ID), events.ChatMessage{
			MessageID: message.ID,
			Sender:    *user,
			Content:   message.Content,
			Type:      message.Type,
			DataJSON:  msg.DataJSON,
		})
	}

	return nil
}

func getCachedUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("chat-user.%v", userID)

	err := database.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return saveDBUser(userID, key)
	}

	userStr, err := database.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return saveDBUser(userID, key)
	}

	var user models.User
	err = json.NewDecoder(strings.NewReader(userStr)).Decode(&user)
	if err != nil {
		return saveDBUser(userID, key)
	}

	return &user, nil
}

func saveDBUser(userID uint, key string) (*models.User, error) {
	user := repository.User.FindById(userID)
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	database.RedisClient.Set(ctx, key, data, time.Minute*20)
	return user, nil
}

func getCachedGroupUserIDs(groupID string) []uint {
	key := fmt.Sprintf("chat-group.%v.user-ids", groupID)
	err := database.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return saveDBGroupUsers(groupID, key)
	}

	gIDsStr, err := database.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return saveDBGroupUsers(groupID, key)
	}

	var gIDs []uint
	err = json.NewDecoder(strings.NewReader(gIDsStr)).Decode(&gIDs)
	if err != nil {
		return saveDBGroupUsers(groupID, key)
	}

	return gIDs
}

func saveDBGroupUsers(groupID string, key string) []uint {
	var gIDs []uint

	gID, err := strconv.Atoi(groupID)
	if err != nil {
		return gIDs
	}

	gIDs = repository.Group.GetUserIDs(uint(gID))
	data, err := json.Marshal(gIDs)
	if err != nil {
		return gIDs
	}

	database.RedisClient.Set(ctx, key, data, time.Minute*20)
	return gIDs
}
