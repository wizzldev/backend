package paginator

import (
	"fmt"
	"gorm.io/gorm"
)

type Config struct {
	Cursor     string
	Order      string
	PointsNext bool
	Limit      int
}

func Paginate[M interface{}](query *gorm.DB, config *Config) (*[]M, string, string, error) {
	var models []M

	isFirstPage := config.Cursor == ""
	pointsNext := false

	if config.Cursor != "" {
		decodedCursor, err := decodeCursor(config.Cursor)
		if err != nil {
			return nil, "", "", err
		}

		pointsNext = decodedCursor["points_next"] == true
		operator, order := getPaginationOperator(pointsNext, config.Order)
		whereStr := fmt.Sprintf("(created_at %s ? or (created_at = ? and id %s ?))", operator, operator)
		query = query.Where(whereStr, decodedCursor["created_at"], decodedCursor["created_at"], decodedCursor["id"])
		if order != "" {
			config.Order = order
		}
	}

	query.Order("created_at " + config.Order).Limit(config.Limit + 1).Find(&models)
	hasPagination := len(models) > config.Limit

	if hasPagination {
		models = models[:config.Limit]
	}

	if !isFirstPage && !pointsNext {
		models = reverse(models)
	}

	next, prev := calculatePagination[M](isFirstPage, hasPagination, config.Limit, models, pointsNext)

	return &models, next, prev, nil
}
