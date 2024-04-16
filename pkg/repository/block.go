package repository

import "github.com/wizzldev/chat/database/models"

type block struct{}

var Block block

func (block) IsBlocked(userID, blockedUserID uint) bool {
	return IsExists[models.Block]([]string{"user_id", "blocked_user_id"}, []any{userID, blockedUserID})
}
