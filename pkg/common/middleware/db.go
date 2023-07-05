package middleware

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/constant"
	"wan_go/sdk"
)

func WithContextDb(c *gin.Context) {
	c.Set(constant.DB, sdk.Runtime.GetDbByKey(c.Request.Host).WithContext(c))
	c.Next()
}
