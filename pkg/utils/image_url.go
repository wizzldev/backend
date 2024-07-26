package utils

import "fmt"

func GetAvatarURL(image string) string {
	return fmt.Sprintf("https://api.wizzl.app/storage/avatars/%s", image)
}
