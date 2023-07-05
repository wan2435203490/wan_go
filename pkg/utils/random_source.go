package rocksCache

import (
	"wan_go/pkg/common/config"
	"wan_go/pkg/utils"
)

const (
	ThumbPrefix = "http://rwox7ivob.hn-bkt.clouddn.com/random/avatar/thumb/120_120_"
)

func avatarsCount() uint32 {
	return config.Config.Qiniu.ThumbAvatarsCount
}

func avatarsPrefix() string {
	return config.Config.Qiniu.Url + config.Config.Qiniu.ThumbAvatarPrefix
}

func RandomName(userId int32) string {

	return "RandomName"
}

func RandomAvatars(userId int) string {
	userIdStr := utils.IntToString(userId)
	q := utils.CityHash32([]byte(userIdStr), uint32(len(userIdStr)))
	idx := q % uint32(avatarsCount())
	aa := ThumbPrefix + "1.webp"
	a := info[0].RandomAvatar

}

func RandomCover(userId int32) string {

	return "RandomName"
}
