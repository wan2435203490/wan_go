package blog

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/utils"
)

// GetUpToken 获取覆盖凭证
func GetUpToken(c *gin.Context) {
	var key string
	a.StringFailed(&key, "key")

	token := utils.GetQiniuToken(key)
	a.OK(token)
}
