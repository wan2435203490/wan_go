package utils

import (
	"fmt"
	"wan_go/pkg/common/config"
)

func avatarsCount() int {
	return config.Config.Qiniu.ThumbAvatarsCount
}

func coverCount() int {
	return config.Config.Qiniu.RandomCoverCount
}

func avatarsPrefix() string {
	return config.Config.Qiniu.Url + config.Config.Qiniu.ThumbAvatarPrefix
}

func randomCoverPrefix() string {
	return config.Config.Qiniu.Url + config.Config.Qiniu.RandomCoverPrefix
}

func randomNames() []string {
	return config.Config.Blog.RandomNames
}

func RandomName(userId int32) string {
	names := randomNames()
	idx := CityHash32Range(userId, 0, len(names)-1)
	return names[idx]
}

func RandomAvatar(userId int32) string {
	idx := CityHash32Range(userId, 1, avatarsCount())
	url := fmt.Sprintf(avatarsPrefix(), idx)
	return url
}

func RandomCover(userId int32) string {
	idx := CityHash32Range(userId, 1, coverCount())
	url := fmt.Sprintf(randomCoverPrefix(), idx)
	return url
}
