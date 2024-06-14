package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/repository"
	"strconv"
	"strings"
	"time"
)

type WSCache struct{}

func NewWSCache() *WSCache {
	return &WSCache{}
}

var ctx = context.Background()

func (w *WSCache) GetUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("chat-user.%v", userID)

	err := rdb.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return w.saveUser(userID, key)
	}

	userStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return w.saveUser(userID, key)
	}

	var user models.User
	err = json.NewDecoder(strings.NewReader(userStr)).Decode(&user)
	if err != nil {
		return w.saveUser(userID, key)
	}

	return &user, nil
}

func (w *WSCache) GetGroupMemberIDs(groupID string) []uint {
	key := fmt.Sprintf("chat-group.%v.user-ids", groupID)
	err := rdb.RedisClient.Exists(ctx, key).Err()
	if err != nil {
		return w.saveGroupMemberIDs(groupID, key)
	}

	gIDsStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return w.saveGroupMemberIDs(groupID, key)
	}

	var gIDs []uint
	err = json.NewDecoder(strings.NewReader(gIDsStr)).Decode(&gIDs)
	if err != nil {
		return w.saveGroupMemberIDs(groupID, key)
	}

	return gIDs
}

func (*WSCache) saveGroupMemberIDs(groupID string, key string) []uint {
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

func (*WSCache) saveUser(userID uint, key string) (*models.User, error) {
	user := repository.User.FindById(userID)
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	rdb.RedisClient.Set(ctx, key, data, time.Minute*20)
	return user, nil
}
