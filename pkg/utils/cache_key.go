package utils

import (
	"fmt"
	"wan_go/pkg/common/constant/blog_const"
)

func CaptchaKey(userId string, flag int) string {
	return fmt.Sprintf(`%s%s_%d`, blog_const.USER_CODE, userId, flag)
}

func CaptchaBindKey(userId, place string, flag int) string {
	return fmt.Sprintf(`%s%s_%s_%d`, blog_const.USER_CODE, userId, place, flag)
}
