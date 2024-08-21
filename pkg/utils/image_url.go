package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetAvatarURL(image string, size ...int) string {
	var img string
	if len(size) == 0 {
		img = image
	} else {
		data := strings.Split(image, ".")
		img = data[0] + "-s" + strconv.Itoa(size[0]) + "." + data[1]
	}
	return fmt.Sprintf("https://api.wizzl.app/storage/avatars/%s", img)
}
