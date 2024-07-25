package utils

import (
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

func IsEmoji(s string) bool {
	_, err := emoji.LookupEmoji(s)
	return err == nil
}
