package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

func MessageActionHandler(s *ws.Server, conn *ws.Connection, userID uint, msg *ws.ClientMessage) error {
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

	switch msg.Type {
	case "message":
		return events.DispatchMessage(s.ID, getCachedGroupUserIDs(s.ID), uint(gID), user, msg)
	case "message.like":
		return events.DispatchMessageLike(s.ID, getCachedGroupUserIDs(s.ID), uint(gID), user, msg)
	default:
		conn.Send(ws.Message{
			Event: "error",
			Data:  fmt.Sprintf("Unknown message type: %s", msg.Type),
		})
	}

	return nil
}

func getCachedUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("chat-user.%v", userID)

	err := rdb.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return saveDBUser(userID, key)
	}

	userStr, err := rdb.RedisClient.Get(ctx, key).Result()
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
	rdb.RedisClient.Set(ctx, key, data, time.Minute*20)
	return user, nil
}

func getCachedGroupUserIDs(groupID string) []uint {
	key := fmt.Sprintf("chat-group.%v.user-ids", groupID)
	err := rdb.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return saveDBGroupUsers(groupID, key)
	}

	gIDsStr, err := rdb.RedisClient.Get(ctx, key).Result()
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
	var uIDs []uint

	gID, err := strconv.Atoi(groupID)
	if err != nil {
		return uIDs
	}

	uIDs = repository.Group.GetUserIDs(uint(gID))
	data, err := json.Marshal(uIDs)
	if err != nil {
		return uIDs
	}

	rdb.RedisClient.Set(ctx, key, data, time.Minute*20)
	return uIDs
}
