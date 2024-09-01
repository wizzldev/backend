package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository/paginator"
)

type theme struct{}

var Theme theme

func (theme) Find(id uint) *models.Theme {
	return FindModelBy[models.Theme]([]string{"id"}, []any{id})
}

func (theme) Paginate(cursor string) (Pagination[models.Theme], error) {
	query := database.DB.Model(&models.Theme{})

	data, next, prev, err := paginator.Paginate[models.Theme](query, &paginator.Config{
		Cursor:     cursor,
		Order:      "desc",
		Limit:      30,
		PointsNext: false,
	})

	return Pagination[models.Theme]{
		Data:       data,
		NextCursor: next,
		Previous:   prev,
	}, err
}
