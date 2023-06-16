package blog

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	blogVO "wan_go/pkg/vo/blog"
)

func ListUser(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blog.User]
	if a.BindFailed(&vo) {
		return
	}

	db_user.ListUser(&vo)

	a.OK(&vo)
}

// ChangeUserStatus
// 修改用户状态
// flag = true：解禁
// flag = false：封禁
func ChangeUserStatus(c *gin.Context) {
	var vo blogVO.BaseRequestVO[*blog.User]
	if a.BindFailed(&vo) {
		return
	}

	db_user.ListUser(&vo)

	a.OK(&vo)
}
