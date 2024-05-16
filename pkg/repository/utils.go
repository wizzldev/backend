package repository

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/database"
	"strings"
)

func FindModelBy[M any](fields []string, values []any) *M {
	var model M

	_ = database.DB.Model(model).
		Where(buildWhereQuery(fields), values...).
		Limit(1).
		Find(&model)

	return &model
}

func All[M any]() []*M {
	var models []*M
	_ = database.DB.Find(&models)
	return models
}

func buildWhereQuery(fields []string) string {
	var query string
	for _, f := range fields {
		query += " " + f + " = ? and"
	}
	return strings.TrimSuffix(query, " and")
}

func IsExists[M any](fields []string, values []any) bool {
	var model M
	var count int64

	err := database.DB.Model(model).
		Where(buildWhereQuery(fields), values...).
		Limit(1).
		Count(&count).
		Error

	if err != nil {
		log.Warn(err)
		return false
	}

	return count > 0
}

func IDsExists[M any](IDs []uint) []uint {
	var model M
	var existing []uint
	database.DB.Model(model).Select("id").Where("id in (?)", IDs).Find(&existing)
	return existing
}
