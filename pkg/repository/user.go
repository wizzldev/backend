package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type user struct{}

var User user

// FindById returns the user with the given id or an error
func (u user) FindById(id uint) *models.User {
	return u.FindBy([]string{"id"}, id)
}

func (u user) FindByEmail(email string) *models.User {
	return u.FindBy([]string{"email"}, email)
}

func (user) FindBy(fields []string, values ...interface{}) *models.User {
	return FindModelBy[models.User](fields, values)
}

func (user) IsEmailExists(email string) bool {
	return IsExists[models.User]([]string{"email"}, []any{email})
}

func (user) IsIPAllowed(uID uint, ip string) bool {
	var count int64
	database.DB.Model(&models.AllowedIP{}).Where("user_id = ? and ip = ? and active = ?", uID, ip, true).Count(&count)
	return count > 0
}

func (user) IsBlocked(blockerID uint, blockedID uint) bool {
	return IsExists[models.Block]([]string{"user_id", "blocked_user_id"}, []any{blockerID, blockedID})
}

func (user) Search(f string, l string, e string, page int) []*models.User {
	var users []*models.User
	q := database.DB.Model(models.User{})
	where := ""
	var whereData []any
	if f != "" {
		where += `first_name like ?`
		whereData = append(whereData, "%"+f+"%")
	}
	if l != "" {
		if where != "" {
			where += " and "
		}
		where += `last_name like ?`
		whereData = append(whereData, "%"+l+"%")
	}

	if e != "" {
		if where != "" {
			where += " and "
		}
		where += `email like ?`
		whereData = append(whereData, "%"+e+"%")
	}

	if where != "" {
		q.Where(where, whereData...)
	}

	err := q.Order("created_at desc").
		Limit(10).
		Offset(10 * (page - 1)).
		Find(&users).Error

	if err != nil {
		return users
	}

	return users
}
