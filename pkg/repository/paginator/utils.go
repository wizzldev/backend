package paginator

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"time"
)

func reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type Cursor map[string]interface{}

func createCursor(id uint, createdAt time.Time, pointsNext bool) Cursor {
	return Cursor{
		"id":          id,
		"created_at":  createdAt,
		"points_next": pointsNext,
	}
}

func encodeCursor(cursor Cursor) string {
	if len(cursor) == 0 {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}
	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

func decodeCursor(cursor string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}
	return cur, nil
}

func getPaginationOperator(pointsNext bool, sortOrder string) (string, string) {
	if pointsNext && sortOrder == "asc" {
		return ">", ""
	}
	if pointsNext && sortOrder == "desc" {
		return "<", ""
	}
	if !pointsNext && sortOrder == "asc" {
		return "<", "desc"
	}
	if !pointsNext && sortOrder == "desc" {
		return ">", "asc"
	}

	return "", ""
}

func calculatePagination[M interface{}](isFirstPage bool, hasPagination bool, limit int, messages []M, pointsNext bool) (string, string) {
	nextCur := Cursor{}
	prevCur := Cursor{}
	if isFirstPage {
		if hasPagination {
			nextCur = createCursor(getAttr(messages[limit-1], "ID").(uint), getAttr(messages[limit-1], "CreatedAt").(time.Time), true)
		}
	} else {
		if pointsNext {
			// if pointing next, it always has prev, but it might not have next
			if hasPagination {
				nextCur = createCursor(getAttr(messages[limit-1], "ID").(uint), getAttr(messages[limit-1], "CreatedAt").(time.Time), true)
			}
			prevCur = createCursor(getAttr(messages[0], "ID").(uint), getAttr(messages[0], "CreatedAt").(time.Time), false)
		} else {
			// this is case of prev, there will always be, nest, but prev needs to be calculated
			nextCur = createCursor(getAttr(messages[limit-1], "ID").(uint), getAttr(messages[limit-1], "CreatedAt").(time.Time), true)
			if hasPagination {
				prevCur = createCursor(getAttr(messages[0], "ID").(uint), getAttr(messages[limit-1], "CreatedAt").(time.Time), false)
			}
		}
	}

	return encodeCursor(nextCur), encodeCursor(prevCur)
}

func getAttr(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}
