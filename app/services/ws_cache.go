package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
	"strconv"
	"time"
)

var ctx = context.Background()

type WSCache struct{}

func NewWSCache() *WSCache {
	return new(WSCache)
}

func (w WSCache) GetUser(userID uint) (*models.User, error) {
	key := w.key(fmt.Sprintf("user:%d", userID))

	userStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return w.getAndSaveUser(userID, key)
	}

	var user models.User

	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return w.getAndSaveUser(userID, key)
	}

	return &user, nil
}

func (WSCache) getAndSaveUser(userID uint, key string) (*models.User, error) {
	user := repository.User.FindById(userID)
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	rdb.RedisClient.Set(ctx, key, data, time.Minute*20)
	return user, nil
}

func (w WSCache) GetGroupMemberIDs(groupID string) []uint {
	key := w.key(fmt.Sprintf("group.%v.userIds", groupID))

	gIDsStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return w.getAndSaveGroupIDs(groupID, key)
	}

	var gIDs []uint
	if err := json.Unmarshal([]byte(gIDsStr), &gIDs); err != nil {
		return w.getAndSaveGroupIDs(groupID, key)
	}

	return gIDs
}

func (WSCache) getAndSaveGroupIDs(groupID string, key string) []uint {
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

func (w WSCache) GetRoles(userID uint, groupID uint) role.Roles {
	key := w.key(fmt.Sprintf("roles.user:%d.%d", userID, groupID))

	roleStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return w.getAndSaveUserRoles(userID, groupID, key)
	}

	var roles []string
	if err = json.Unmarshal([]byte(roleStr), &roles); err != nil {
		return w.getAndSaveUserRoles(userID, groupID, key)
	}

	return *role.NewRoles(roles)
}

func (w WSCache) getAndSaveUserRoles(userID uint, groupID uint, key string) role.Roles {
	roles := repository.Group.GetUserRoles(groupID, userID, w.GetGroupRoles(groupID))
	rdb.RedisClient.Set(ctx, key, roles.String(), time.Minute*20)
	return roles
}

func (w WSCache) GetGroupRoles(groupID uint) role.Roles {
	key := w.key(fmt.Sprintf("roles.group:%d", groupID))

	gIDsStr, err := rdb.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return *w.getAndSaveGroupRoles(groupID, key)
	}

	var roles []string
	if err := json.Unmarshal([]byte(gIDsStr), &roles); err != nil {
		return *w.getAndSaveGroupRoles(groupID, key)
	}

	return *role.NewRoles(roles)
}

func (WSCache) getAndSaveGroupRoles(groupID uint, key string) *role.Roles {
	group := repository.Group.Find(groupID)
	if group.ID < 1 {
		return new(role.Roles)
	}

	roles := role.NewRoles(group.Roles)

	rdb.RedisClient.Set(ctx, key, roles.String(), time.Minute*20)

	return roles
}

func (w WSCache) IsPM(groupID uint) bool {
	// TODO: make it cacheable
	return repository.Group.Find(groupID).IsPrivateMessage
}

func (WSCache) key(s string) string {
	return "ws-" + s
}
